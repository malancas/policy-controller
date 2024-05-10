package customauth

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	kauth "github.com/google/go-containerregistry/pkg/authn/kubernetes"
	"k8s.io/client-go/kubernetes"
)

func K8sChainWithCustomACR(ctx context.Context, client kubernetes.Interface, opt k8schain.Options) (authn.Keychain, error) {
	k8s, err := kauth.New(ctx, client, kauth.Options(opt))
	if err != nil {
		return nil, err
	}

	return authn.NewMultiKeychain(
		k8s,
		authn.DefaultKeychain,
		authn.NewKeychainFromHelper(NewACRHelper()),
	), nil

}

type ACRHelper struct{}

func NewACRHelper() credentials.Helper {
	return &ACRHelper{}
}

func (a ACRHelper) Add(_ *credentials.Credentials) error {
	return fmt.Errorf("add is unimplemented")
}

func (a ACRHelper) Delete(_ string) error {
	return fmt.Errorf("delete is unimplemented")
}

func (a ACRHelper) Get(serverURL string) (string, string, error) {
	azCred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	token, err := azCred.GetToken(context.Background(), policy.TokenRequestOptions{})
	if err != nil {
		log.Fatalf("failed to get token: %v", err)
	}

	return token.Token, "", nil
}

func (a ACRHelper) List() (map[string]string, error) {
	return nil, fmt.Errorf("list is unimplemented")
}
