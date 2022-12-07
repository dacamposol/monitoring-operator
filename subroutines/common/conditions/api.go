package conditions

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate go run -mod=mod github.com/vektra/mockery/v2 --all --case=underscore --with-expecter

type Setter interface {
	SetTrue(objectMeta metav1.ObjectMeta, conditions *[]metav1.Condition, reason, message string)
	SetFalse(objectMeta metav1.ObjectMeta, conditions *[]metav1.Condition, reason, message string)
}
