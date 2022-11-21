package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	authv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var kClientset *kubernetes.Clientset

func arrayContain(t []string, needle string) bool {
	for i := range t {
		if t[i] == needle {
			return true
		}
	}
	return false
}

// https://stackoverflow.com/a/51270134
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func setup() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	kClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}
func verifyToken(clientId string) (*authv1.TokenReview, bool, error) {
	tr := authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token: clientId,
		},
	}
	result, err := kClientset.AuthenticationV1().TokenReviews().Create(context.TODO(), &tr, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}

	log.Printf("%v\n", prettyPrint(result.Status))
	if result.Status.Authenticated {
		return result, true, nil
	}
	return nil, false, nil

}

func getRoleBinding(s *authv1.TokenReview) (*v1.PolicyRule, error) {
	roleBindingList, err := kClientset.RbacV1().RoleBindings(s.Namespace).List(context.TODO(), metav1.ListOptions{
		// use this label to filter
		LabelSelector: "usedby=storage-hub",
	})
	if err != nil {
		log.Println("ERROR get rolebinding= ", err)
		return nil, err
	}
	for _, val := range roleBindingList.Items {
		for _, sub := range val.Subjects {
			serviceAccountName := strings.Split(s.Status.User.Username, ":")[len(strings.Split(s.Status.User.Username, ":"))-1]
			serviceAccountNamespace := strings.Split(s.Status.User.Username, ":")[len(strings.Split(s.Status.User.Username, ":"))-2]
			if sub.Name == serviceAccountName && sub.Namespace == serviceAccountNamespace {
				role, err := kClientset.RbacV1().Roles(serviceAccountNamespace).Get(context.TODO(), val.RoleRef.Name, metav1.GetOptions{})
				if err != nil {
					return nil, err
				}
				for _, rule := range role.Rules {
					if arrayContain(rule.APIGroups, "storage-hub") && arrayContain(rule.Resources, "database") {
						log.Printf("this user is able to '%v', on this resources '%v' on this APIGroups '%v'", rule.Verbs, rule.Resources, rule.APIGroups)
						return rule.DeepCopy(), nil
					}
				}
			}
		}

	}
	return nil, nil
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("NEW REQUEST")
	// Read the value of the client identifier from the request header
	clientId := r.Header.Get("Bearer")
	if len(clientId) == 0 {
		http.Error(w, "Bearer not supplied", http.StatusUnauthorized)
		return
	}
	serviceAccountInfo, authenticated, err := verifyToken(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !authenticated {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}
	// // cerise sur le gateau => check IP
	// ip := r.Header.Get("X-Real-IP")
	// if len(ip) == 0 {
	// 	http.Error(w, "X-Real-IP not supplied", http.StatusUnauthorized)
	// 	return
	// }
	// if err := ceriseSurLeGateau(serviceAccountInfo, ip); err != nil {
	// 	http.Error(w, err.Error(), http.StatusForbidden)
	// 	return
	// }
	log.Println("authenticated user authenticated")
	// Get role permission
	rule, err := getRoleBinding(serviceAccountInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rule == nil {
		http.Error(w, "you have no permission on storage-hub", http.StatusForbidden)
		return
	}
	// OK
	io.WriteString(w, fmt.Sprintf("Hello from storage-hub. You have been authenticated and you have the permission to '%v' on this resources '%v'. Enjoy!", rule.Verbs, rule.Resources))

}

func main() {
	setup()

	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8081", nil)
}

// func ceriseSurLeGateau(s *authv1.TokenReview, ip string) error {
// 	if len(s.Status.User.Extra["authentication.kubernetes.io/pod-name"]) == 0 {
// 		log.Println("call from outside")
// 		return nil
// 	}
// 	serviceAccountNamespace := strings.Split(s.Status.User.Username, ":")[len(strings.Split(s.Status.User.Username, ":"))-2]
// 	pod, err := kClientset.CoreV1().Pods(serviceAccountNamespace).Get(context.TODO(), s.Status.User.Extra["authentication.kubernetes.io/pod-name"][0], metav1.GetOptions{})
// 	if err != nil {
// 		log.Println("ERROR get pod= ", err)
// 		return err
// 	}

// 	if pod.Status.PodIP != ip {
// 		log.Println("ip expected =", pod.Status.PodIP)
// 		log.Println("ip given =", ip)
// 		return errors.New("you are trying to use a token from an other client")
// 	}
// 	return nil
// }
