package usecase_test

import (
	"bytes"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/backup-blob/zfs-backup-blob/internal/domain/mocks"
	"github.com/backup-blob/zfs-backup-blob/internal/usecase"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"io"
	"testing"
)

func TestSpecVolume(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockVolumeRepo := mocks.NewMockVolumeRepository(ctrl)
	mockRenderRepo := mocks.NewMockRenderRepository(ctrl)
	mockSnapUsecase := mocks.NewMockSnapshotUsecase(ctrl)
	volumes := []*domain.ZfsVolume{
		{Name: "/vol/vol1", GroupName: "-"},
		{Name: "/vol/vol2", GroupName: "group1"},
	}
	volumeUsecase := usecase.NewVolumeUsecase(mockVolumeRepo, mockRenderRepo, mockSnapUsecase)

	Convey("Given the AddToGroup function is called", t, func() {
		Convey("When everything goes positive", func() {
			Convey("It should add volume to group", func() {
				mockVolumeRepo.EXPECT().TagVolumeWithGroup(gomock.Any()).DoAndReturn(func(v *domain.ZfsVolume) error {
					So(v.Name, ShouldEqual, "pool/path")
					So(v.GroupName, ShouldEqual, "group1")

					return nil
				})
				mockSnapUsecase.EXPECT().CreateByVolume(gomock.Any(), domain.Full).Return(nil)
				err := volumeUsecase.AddToGroup("pool/path", "group1")

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given the List function is called", t, func() {
		Convey("When everything works as expected", func() {
			Convey("It should render a list of volumes", func() {
				writer := &bytes.Buffer{}
				var volumesRendered []string
				mockVolumeRepo.EXPECT().ListVolumes().Return(volumes, nil)
				mockRenderRepo.EXPECT().RenderVolumeTable(writer, gomock.Any()).DoAndReturn(func(w io.Writer, v []*domain.ZfsVolume) {
					for _, vol := range v {
						volumesRendered = append(volumesRendered, vol.Name)
					}
				})

				err := volumeUsecase.List(writer)

				So(err, ShouldBeNil)
				So(volumesRendered, ShouldContain, "/vol/vol2")
				So(volumesRendered, ShouldHaveLength, 1)
			})
		})
	})
}
