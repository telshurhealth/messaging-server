package app

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

func (api *PluginAPI) BatchAddChannelMember(channelIds []string, userID string) ([]*model.ChannelMember, *model.AppError) {
	if len(channelIds) == 0 {
		return []*model.ChannelMember{}, nil
	}

	channels, err := api.app.GetChannelsById(channelIds)
	if err != nil {
		return nil, err
	}

	return api.app.BatchAddChannelMember(api.ctx, userID, channels)
}
