package auth

import (
	"fmt"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"spyrosmoux/api/internal/helpers"
)

var (
	clientId     = helpers.LoadEnvVariable("CLIENT_ID")
	clientSecret = helpers.LoadEnvVariable("CLIENT_SECRET")
)

func InitSuperTokens() {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:3567",
			APIKey:        "someKey", // TODO(spyromoux) use a secure string
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "FlowForge",
			APIDomain:       "http://localhost:8080", // TODO(spyrosmoux) use env variables
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
				Override: &tpmodels.OverrideStruct{
					Functions: overrideThirdPartySignInUp,
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

func overrideThirdPartySignInUp(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
	// create a copy of the originalImplementation
	originalSignInUp := *originalImplementation.SignInUp

	// override the sign in up function
	*originalImplementation.SignInUp = func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens map[string]interface{}, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext *map[string]interface{}) (tpmodels.SignInUpResponse, error) {
		// First we call the original implementation of SignInUp.
		response, err := originalSignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
		if err != nil {
			return tpmodels.SignInUpResponse{}, err
		}

		fmt.Println(response.OK.OAuthTokens["access_token"].(string))

		if response.OK != nil {
			// sign in / up was successful
			if response.OK.CreatedNewUser {
				// TODO(spyrosmoux) save the user id in api db
			}
		}
		return response, nil
	}

	return originalImplementation
}
