package middlewares

import (
	"github.com/spyrosmoux/cicd/common/helpers"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	clientId     = helpers.LoadEnvVariable("GITHUB_OAUTH_APP_CLIENT_ID")
	clientSecret = helpers.LoadEnvVariable("GITHUB_OAUTH_APP_CLIENT_SECRET")
)

func InitSuperTokens() {
	apiBasePath := "/app/cicd/api/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:3567",
			APIKey:        "someKey", // TODO(spyromoux) use a secure string
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "FlowForge",
			APIDomain:       "http://localhost:8080/api", // TODO(@SpyrosMoux) use env variables
			WebsiteDomain:   "http://localhost:3000",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "github",
								Clients: []tpmodels.ProviderClientConfig{
									{
										// Refers to OAuth GitHub App
										ClientID:     clientId,
										ClientSecret: clientSecret,
									},
								},
							},
						},
					},
				},
			}),
			emailpassword.Init(nil),
			session.Init(&sessmodels.TypeInput{
				ExposeAccessTokenToFrontendInCookieBasedAuth: true,
			}), // initializes session features
		},
	})

	if err != nil {
		panic(err.Error())
	}
}
