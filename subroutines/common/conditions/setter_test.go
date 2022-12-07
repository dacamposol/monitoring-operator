package conditions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConditionsSetter(t *testing.T) {
	var conditions []metav1.Condition

	t.Run("New Conditions Setter", func(t *testing.T) {
		// Arrange + Act
		setter := NewSetter("MyConditionType")

		// Assert
		assert.NotNil(t, setter)
	})

	t.Run("Set to True", func(t *testing.T) {
		// Arrange
		setter := NewSetter("MyConditionType")

		// Act
		setter.SetTrue(metav1.ObjectMeta{}, &conditions, "reason", "message")

		// Assert
		assert.Equal(t, len(conditions), 1)
		assert.Equal(t, conditions[0].Status, metav1.ConditionTrue)
		assert.Equal(t, conditions[0].Type, "MyConditionType")
		assert.Equal(t, conditions[0].ObservedGeneration, int64(0))
		assert.Equal(t, conditions[0].Reason, "reason")
		assert.Equal(t, conditions[0].Message, "message")
	})

	t.Run("Set to False", func(t *testing.T) {
		// Arrange
		setter := NewSetter("MyConditionType")

		// Act
		setter.SetFalse(metav1.ObjectMeta{}, &conditions, "reason", "message")

		// Assert
		assert.Equal(t, len(conditions), 1)
		assert.Equal(t, conditions[0].Status, metav1.ConditionFalse)
		assert.Equal(t, conditions[0].Type, "MyConditionType")
		assert.Equal(t, conditions[0].ObservedGeneration, int64(0))
		assert.Equal(t, conditions[0].Reason, "reason")
		assert.Equal(t, conditions[0].Message, "message")
	})
}
