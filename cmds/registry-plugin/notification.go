package main

import (
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/ppacher/portmaster-plugin-registry/manager"
	"github.com/ppacher/portmaster-plugin-registry/structs"
	"github.com/safing/portmaster/plugin/framework"
	"github.com/safing/portmaster/plugin/shared/notification"
	"github.com/safing/portmaster/plugin/shared/proto"
)

type NotificationHandler struct {
	notification.Service

	manager *manager.Manager
}

func NewNotificationHandler(manager *manager.Manager, notify notification.Service) *NotificationHandler {
	handler := &NotificationHandler{
		Service: notify,
		manager: manager,
	}

	manager.OnFetchDone(handler.onFetchDone)
	manager.OnUpdateAvailable(handler.onUpdateAvailable)

	_, err := framework.Notify().CreateNotification(framework.Context(), &proto.Notification{
		EventId:      "plugin-registry:peristent-notification",
		Title:        "PECS",
		Message:      "PECS is successfully installed and running",
		ShowOnSystem: false,
		Actions: []*proto.NotificationAction{
			{
				Id:   "open",
				Text: "Open",
				ActionType: &proto.NotificationAction_OpenUrl{
					OpenUrl: "https://github.com/ppacher/portmaster-plugin-registry",
				},
			},
		},
	})
	if err != nil {
		hclog.L().Error("failed to create persistant info notification", "error", err.Error())
	}

	return handler
}

func (handler *NotificationHandler) onFetchDone(err error) {
	if err != nil {
		_, err := handler.CreateNotification(framework.Context(), &proto.Notification{
			EventId: "plugin-registry:fetch-failed",
			Type:    proto.NotificationType_NOTIFICATION_TYPE_ERROR,
			Title:   "Failed to fetch plugin repositories",
			Message: err.Error(),
			Actions: []*proto.NotificationAction{
				{
					Id:   "go-away",
					Text: "OK",
				},
			},
		})
		if err != nil {
			hclog.L().Error("failed to create fetch-failed notification", "error", err)
		}
	} else {
		_, err := handler.CreateNotification(framework.Context(), &proto.Notification{
			EventId: "plugin-registry:fetch-failed",
			Type:    proto.NotificationType_NOTIFICATION_TYPE_ERROR,
			Expires: time.Now().Add(-time.Second).UnixNano(),
		})
		if err != nil {
			hclog.L().Error("failed to clear fetch-failed notification", "error", err)
		}
	}
}

func (handler *NotificationHandler) onUpdateAvailable(updates []structs.AvailableUpdate) {
	for _, upd := range updates {
		_, err := handler.CreateNotification(framework.Context(), &proto.Notification{
			EventId: "plugin-registry:update-" + upd.Name,
			Type:    proto.NotificationType_NOTIFICATION_TYPE_INFO,
			Title:   upd.Name + ": new version " + upd.NewVersion + " is available",
			Message: "A new version for the plugin " + upd.Name + " is available. Update to " + upd.NewVersion + " now?",
			Actions: []*proto.NotificationAction{
				{
					Id:   "update-now",
					Text: "Update Now",
				},
				{
					Id:   "not-now",
					Text: "Later",
				},
			},
		})

		if err != nil {
			hclog.L().Error("failed to create update notification", "plugin", upd.Name, "error", err.Error())
		}
	}
}
