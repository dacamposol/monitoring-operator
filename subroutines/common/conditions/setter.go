package conditions

import (
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type conditionsSetter struct {
	conditionType string
}

func NewSetter(conditionType string) Setter {
	return &conditionsSetter{conditionType: conditionType}
}

// Functions

func setCondition(objectMeta metav1.ObjectMeta, conditions *[]metav1.Condition, conditionType string, status metav1.ConditionStatus, reason, message string) {
	workflowCondition := metav1.Condition{
		Type:               conditionType,
		Status:             status,
		ObservedGeneration: objectMeta.GetGeneration(),
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             reason,
		Message:            message,
	}
	meta.SetStatusCondition(conditions, workflowCondition)
}

func (c *conditionsSetter) SetTrue(objectMeta metav1.ObjectMeta, conditions *[]metav1.Condition, reason, message string) {
	setCondition(objectMeta, conditions, c.conditionType, metav1.ConditionTrue, reason, message)
}

func (c *conditionsSetter) SetFalse(objectMeta metav1.ObjectMeta, conditions *[]metav1.Condition, reason, message string) {
	setCondition(objectMeta, conditions, c.conditionType, metav1.ConditionFalse, reason, message)
}
