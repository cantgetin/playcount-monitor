package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"playcount-monitor-backend/internal/database/repository"
	"playcount-monitor-backend/internal/database/repository/model"
	"playcount-monitor-backend/internal/dto"
	usercardcreate "playcount-monitor-backend/internal/usecase/usercard/create"
	usercardupdate "playcount-monitor-backend/internal/usecase/usercard/update"
	"playcount-monitor-backend/tests/integration"
	"time"
)

func (s *IntegrationSuite) Test_CreateUseCard() {
	s.Run("valid requests", func() {
		tt := []struct {
			name    string
			in      *usercardcreate.CreateUserCardCommand
			outCode int
		}{
			{
				name: "valid request, should properly create",
				in: &usercardcreate.CreateUserCardCommand{
					User: &dto.CreateUserCommand{
						ID:                       1,
						AvatarURL:                "avatarurl.com",
						Username:                 "username",
						UnrankedBeatmapsetCount:  1,
						GraveyardBeatmapsetCount: 1,
					},
					Mapsets: []*dto.CreateMapsetCommand{
						{
							Id:     1,
							Artist: "artist",
							Title:  "title",
							Covers: map[string]string{
								"cover1": "cover1",
								"cover2": "cover2",
							},
							Status:         "graveyard",
							LastUpdated:    time.Now().UTC(),
							UserId:         1,
							PreviewUrl:     "previewurl.com",
							Tags:           "tags tags",
							PlayCount:      20,
							FavouriteCount: 25,
							Bpm:            150,
							Creator:        "username",
							Beatmaps: []*dto.CreateBeatmapCommand{
								{
									Id:               1,
									BeatmapsetId:     1,
									DifficultyRating: 5.3,
									Version:          "version1",
									Accuracy:         6.7,
									Ar:               9.3,
									Bpm:              150.3,
									Cs:               4,
									Status:           "graveyard",
									Url:              "beatmapurl.com",
									TotalLength:      3,
									UserId:           1,
									Passcount:        12,
									Playcount:        13,
									LastUpdated:      time.Now().UTC(),
								},
								{
									Id:               2,
									BeatmapsetId:     1,
									DifficultyRating: 6.8,
									Version:          "version2",
									Accuracy:         4.6,
									Ar:               9,
									Bpm:              150,
									Cs:               5,
									Status:           "graveyard",
									Url:              "beatmap2url.com",
									TotalLength:      4,
									UserId:           1,
									Passcount:        0,
									Playcount:        7,
									LastUpdated:      time.Now().UTC(),
								},
							},
						},
					},
				},
				outCode: 200,
			},
		}

		for _, tc := range tt {
			s.Run(tc.name, func() {
				inJSON, err := json.Marshal(tc.in)
				s.Require().NoError(err)

				out, err := http.Post(
					"http://localhost:8080/user_card/create",
					"application/json",
					bytes.NewBuffer(inJSON),
				)

				s.Require().NoError(err)
				s.Require().NotNil(out)

				s.Assert().Equal(out.StatusCode, tc.outCode)
			})
		}

		err := integration.ClearTables(s.ctx, s.db)
		s.Require().NoError(err)
	})
}

// what if user got new mapset which is not created, handle this case
func (s *IntegrationSuite) Test_UpdateUserCard() {
	s.Run("valid requests", func() {
		type models struct {
			User     *model.User
			Mapset   []*model.Mapset
			Beatmaps []*model.Beatmap
		}

		var tt = []struct {
			name    string
			create  *models
			in      *usercardupdate.UpdateUserCardCommand
			result  *models // assert db models after calling update method
			outCode int
		}{
			{
				name: "valid request, should properly update",
				create: &models{
					User: &model.User{
						ID:                       123,
						Username:                 "username1",
						AvatarURL:                "avararurl.com",
						GraveyardBeatmapsetCount: 1,
						UnrankedBeatmapsetCount:  1,
						CreatedAt:                time.Now().UTC(),
						UpdatedAt:                time.Now().UTC(),
					},
					Mapset: []*model.Mapset{
						{
							ID:          123,
							Artist:      "artist",
							Title:       "title",
							Covers:      repository.JSON(`{"cover1": "cover1", "cover2": "cover2"}`),
							Status:      "graveyard",
							LastUpdated: time.Now().UTC(),
							UserID:      123,
							Creator:     "username1",
							PreviewURL:  "avararurl.com",
							Tags:        "tags tags",
							BPM:         150,
							MapsetStats: repository.JSON(`{"2023-12-24T12:00:00Z":{"play_count":52,"favourite_count":2}}`),
							CreatedAt:   time.Now().UTC(),
							UpdatedAt:   time.Now().UTC(),
						},
					},
					Beatmaps: []*model.Beatmap{
						{
							ID:               3,
							MapsetID:         123,
							DifficultyRating: 5.3,
							Version:          "version1",
							Accuracy:         7.3,
							AR:               9,
							BPM:              150.3,
							CS:               4,
							Status:           "graveyard",
							URL:              "beatmap1url.com",
							TotalLength:      23,
							UserID:           123,
							LastUpdated:      time.Now().UTC(),
							BeatmapStats:     repository.JSON(`{"2023-12-24T12:00:00Z":{"play_count":25,"pass_count":23}}`),
							CreatedAt:        time.Now().UTC(),
							UpdatedAt:        time.Now().UTC(),
						},
						{
							ID:               4,
							MapsetID:         123,
							DifficultyRating: 6.7,
							Version:          "version2",
							Accuracy:         8.5,
							AR:               9.3,
							BPM:              150.2,
							CS:               3.3,
							Status:           "graveyard",
							URL:              "beatmap2url.com",
							TotalLength:      24,
							UserID:           123,
							LastUpdated:      time.Now().UTC(),
							BeatmapStats:     repository.JSON(`{"2023-12-24T12:00:00Z":{"play_count":27,"pass_count":24}}`),
							CreatedAt:        time.Now().UTC(),
							UpdatedAt:        time.Now().UTC(),
						},
					},
				},
				in: &usercardupdate.UpdateUserCardCommand{
					User: &dto.CreateUserCommand{
						ID:                       123,
						AvatarURL:                "avararurlchanged.com",
						Username:                 "username1changed",
						UnrankedBeatmapsetCount:  2, // assume user now have 2 mapsets
						GraveyardBeatmapsetCount: 2,
					},
					Mapsets: []*dto.CreateMapsetCommand{
						{
							Id:     123,
							Artist: "artistchanged",
							Title:  "titlechanged",
							Covers: map[string]string{
								"cover1changed": "cover1changed",
								"cover2changed": "cover2changed",
							},
							Status:         "statuschanged",
							LastUpdated:    time.Now().UTC(),
							UserId:         123,
							PreviewUrl:     "previewurlchanged.com",
							Tags:           "tagschanged tagschanged",
							PlayCount:      200,
							FavouriteCount: 200,
							Bpm:            200,
							Creator:        "username1changed",
							Beatmaps: []*dto.CreateBeatmapCommand{
								{
									Id:               3,
									BeatmapsetId:     123,
									DifficultyRating: 7.6,
									Version:          "version1changed",
									Accuracy:         1,
									Ar:               1,
									Bpm:              1,
									Cs:               1,
									Status:           "graveyard",
									Url:              "urlchanged.com",
									TotalLength:      1,
									UserId:           123,
									Passcount:        100,
									Playcount:        100,
									LastUpdated:      time.Now().UTC(),
								},
								{
									Id:               4,
									BeatmapsetId:     123,
									DifficultyRating: 1.2,
									Version:          "version2changed",
									Accuracy:         2,
									Ar:               2,
									Bpm:              2,
									Cs:               2,
									Status:           "graveyard",
									Url:              "urlchanged.com",
									TotalLength:      1,
									UserId:           123,
									Passcount:        100,
									Playcount:        100,
									LastUpdated:      time.Now().UTC(),
								},
							},
						},
						{
							Id:     345,
							Artist: "artist",
							Title:  "title",
							Covers: map[string]string{
								"cover1": "cover1",
								"cover2": "cover2",
							},
							Status:         "graveyard",
							LastUpdated:    time.Now().UTC(),
							UserId:         123,
							PreviewUrl:     "previewurlnewmap.com",
							Tags:           "tags tags",
							PlayCount:      345,
							FavouriteCount: 456,
							Bpm:            120,
							Creator:        "username1changed",
							Beatmaps: []*dto.CreateBeatmapCommand{
								{
									Id:               1488,
									BeatmapsetId:     345,
									DifficultyRating: 1,
									Version:          "version1",
									Accuracy:         2,
									Ar:               3,
									Bpm:              4,
									Cs:               5,
									Status:           "graveyard",
									Url:              "url.com",
									TotalLength:      1,
									UserId:           2,
									Passcount:        3,
									Playcount:        4,
									LastUpdated:      time.Now().UTC(),
								},
								{
									Id:               1337,
									BeatmapsetId:     345,
									DifficultyRating: 1,
									Version:          "version2",
									Accuracy:         3,
									Ar:               4,
									Bpm:              5,
									Cs:               6,
									Status:           "graveyard",
									Url:              "url.com",
									TotalLength:      1,
									UserId:           2,
									Passcount:        3,
									Playcount:        4,
									LastUpdated:      time.Now().UTC(),
								},
							},
						},
					},
				},
				outCode: 0,
				result: &models{
					User: &model.User{
						ID:                       123,
						Username:                 "username1changed",
						AvatarURL:                "avatarurlchanged.com",
						GraveyardBeatmapsetCount: 2,
						UnrankedBeatmapsetCount:  2,
					},
					Mapset: []*model.Mapset{
						{
							ID:          123,
							Artist:      "artistchanged",
							Title:       "titlechanged",
							Covers:      repository.JSON(``), // todo
							Status:      "statuschanged",
							LastUpdated: time.Now().UTC(),
							UserID:      123,
							Creator:     "username1changed",
							PreviewURL:  "previewurlchanged.com",
							Tags:        "tagschanged tagschanged",
							BPM:         200,
							MapsetStats: nil, // todo
						},
						{
							ID:          345,
							Artist:      "artist",
							Title:       "title",
							Covers:      repository.JSON(``), // todo
							Status:      "graveyard",
							LastUpdated: time.Now().UTC(),
							UserID:      123,
							Creator:     "username1changed",
							PreviewURL:  "previewurlnewmap.com",
							Tags:        "tags tags",
							BPM:         120,
							MapsetStats: nil, // todo
						},
					},
					Beatmaps: nil,
				},
			},
		}
		for _, tc := range tt {
			s.Run(tc.name, func() {
				inJSON, err := json.Marshal(tc.in)
				s.Require().NoError(err)

				out, err := http.Post(
					"http://localhost:8080/user_card/Update",
					"application/json",
					bytes.NewBuffer(inJSON),
				)

				s.Require().NoError(err)
				s.Assert().Equal(out.StatusCode, tc.outCode)
			})
		}

		err := integration.ClearTables(s.ctx, s.db)
		s.Require().NoError(err)
	})
}
