package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/app/request"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/store"
)

func (a *App) GetChannelsById(channelIDs []string) (model.ChannelList, *model.AppError) {
	channels, err := a.Srv().Store.Channel().GetChannelsByIds(channelIDs, true)

	if err != nil {
		return nil, model.NewAppError("GetChannelsById", "app.channel.get_channels_by_ids.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return channels, nil
}

func (a *App) BatchAddChannelMember(c *request.Context, userID string, channels model.ChannelList) ([]*model.ChannelMember, *model.AppError) {
	var user *model.User
	var err *model.AppError

	if user, err = a.GetUser(userID); err != nil {
		return nil, err
	}

	channelMembers := make([]*model.ChannelMember, len(channels))

	for _, channel := range channels {
		if member, err := a.Srv().Store.Channel().GetMember(context.Background(), channel.Id, userID); err != nil {
			var nfErr *store.ErrNotFound
			if !errors.As(err, &nfErr) {
				return nil, model.NewAppError("BatchAddChannelMember", "app.channel.get_member.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		} else {
			channelMembers = append(channelMembers, member)
			continue
		}

		member, err := a.addUserToChannelWithoutMessage(user, channel)
		if err != nil {
			return nil, err
		}

		channelMembers = append(channelMembers, member)
	}

	message := model.NewWebSocketEvent(model.WebsocketEventUserRemoved, channels[0].TeamId, "", user.Id, nil)
	message.Add("user_id", user.Id)
	message.Add("team_id", channels[0].TeamId)
	a.Publish(message)

	return channelMembers, nil
}

func (a *App) addUserToChannelWithoutMessage(user *model.User, channel *model.Channel) (*model.ChannelMember, *model.AppError) {
	newMember, err := a.addUserToChannel(user, channel)
	if err != nil {
		return nil, err
	}
	return newMember, nil
}
