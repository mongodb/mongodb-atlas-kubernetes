package workflow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func Test_EnsureCondition(t *testing.T) {
	t.Run("Adding different conditions", func(t *testing.T) {
		st := &Status{conditions: []status.Condition{}}
		st.EnsureCondition(status.Condition{Type: status.ProjectReadyType})
		st.EnsureCondition(status.Condition{Type: status.IPAccessListReadyType})

		// We cannot check arrays for equality as there's a mutable field LastTransitionTime
		assert.Len(t, st.conditions, 2)
		assert.Equal(t, status.ProjectReadyType, st.conditions[0].Type)
		assert.GreaterOrEqual(t, metav1.Now().Unix(), st.conditions[0].LastTransitionTime.Unix())
		assert.Equal(t, status.IPAccessListReadyType, st.conditions[1].Type)
		assert.GreaterOrEqual(t, metav1.Now().Unix(), st.conditions[1].LastTransitionTime.Unix())
	})
	t.Run("Adding the same conditions the same statuses", func(t *testing.T) {
		st := &Status{conditions: []status.Condition{}}
		st.EnsureCondition(status.Condition{Type: status.IPAccessListReadyType, Status: corev1.ConditionTrue})
		firstCondition := *st.conditions[0].DeepCopy()
		assert.GreaterOrEqual(t, metav1.Now().Unix(), st.conditions[0].LastTransitionTime.Unix())
		assert.Equal(t, status.IPAccessListReadyType, st.conditions[0].Type)
		assert.Equal(t, corev1.ConditionTrue, st.conditions[0].Status)

		time.Sleep(time.Millisecond * 100)
		// We are ensuring the same condition with the same status - the LastTransitionTime should be the same
		st.EnsureCondition(status.Condition{Type: status.IPAccessListReadyType, Status: corev1.ConditionTrue})

		assert.Len(t, st.conditions, 1)
		// Note, that condition is the same after update
		assert.Equal(t, firstCondition, st.conditions[0])
	})
	t.Run("Adding the same conditions different statuses", func(t *testing.T) {
		st := &Status{conditions: []status.Condition{}}
		st.EnsureCondition(status.Condition{Type: status.IPAccessListReadyType, Status: corev1.ConditionTrue})
		firstCondition := *st.conditions[0].DeepCopy()
		assert.GreaterOrEqual(t, metav1.Now().Unix(), st.conditions[0].LastTransitionTime.Unix())
		assert.Equal(t, status.IPAccessListReadyType, st.conditions[0].Type)
		assert.Equal(t, corev1.ConditionTrue, st.conditions[0].Status)

		time.Sleep(time.Millisecond * 100)
		// We are ensuring the same condition with the same status - the LastTransitionTime should be the same
		st.EnsureCondition(status.Condition{Type: status.IPAccessListReadyType, Status: corev1.ConditionFalse})

		assert.Len(t, st.conditions, 1)

		assert.Equal(t, corev1.ConditionFalse, st.conditions[0].Status)
		assert.NotEqual(t, firstCondition.LastTransitionTime, st.conditions[0].LastTransitionTime)
	})
}
