package interfaces

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces/fake"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestListItemWrapper(t *testing.T) {
	appVersionList := v1alpha3.KeptnAppVersionList{
		Items: []v1alpha3.KeptnAppVersion{
			{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StateFailed,
					CurrentPhase: "test",
				},
			},
		},
	}

	object, err := NewListItemWrapperFromClientObjectList(&appVersionList)
	require.Nil(t, err)

	items := object.GetItems()
	require.Len(t, items, 1)
}

func TestListItem(t *testing.T) {
	listItemMock := fake.ListItemMock{
		GetItemsFunc: func() []client.Object {
			return nil
		},
	}
	wrapper := ListItemWrapper{Obj: &listItemMock}
	_ = wrapper.GetItems()
	require.Len(t, listItemMock.GetItemsCalls(), 1)
}
