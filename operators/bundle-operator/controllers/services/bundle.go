package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/bundle-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/bundle-operator/bundles"
	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/verify"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"
)

type BundleService struct {
}

func NewBundleService() *BundleService {
	return &BundleService{}
}

func (bs *BundleService) CheckBundleSignature(ctx context.Context, cr *v1alpha1.EntandoBundleV2, log logr.Logger) (map[string]string, error) {
	m := make(map[string]string, len(cr.Spec.TagList))
	for _, tag := range cr.Spec.TagList {
		if len(tag.SignatureInfo) <= 0 {
			return nil, fmt.Errorf("error signature info empty")
		}
		verifyOk := true
		for _, signature := range tag.SignatureInfo {
			err := bs.verifySignature(cr.Spec.Repository+"@"+tag.Digest, signature.PubKeySecret)
			if err != nil {
				verifyOk = false
				log.Error(err, "error verify signature ",
					"tag", tag.Tag, "digest", tag.Digest, "signType", signature.Type)
			}

		}
		if verifyOk {
			key := utility.TruncateString("signature-"+strings.Split(tag.Digest, ":")[1], 63)
			m[key] = "Verified"
		}
		/*
			ref, err := bs.retrieveSignatureImageRef(cr.Spec.Repository + "@" + tag.Digest)
			if err != nil {
				return nil, err
			}
			fmt.Println("ref signature: " + ref)
		*/
	}
	return m, nil
}

func (bs *BundleService) GenerateBundleCode(cr *v1alpha1.EntandoBundleV2) string {
	s := utility.GenerateSha256(cr.Spec.Repository)
	return "bundle-" + strings.ToLower(utility.TruncateString(s, 8))
}

func (bs *BundleService) GetComponents(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) ([]bundles.Component, string, error) {
	/*
		repository := "docker.io/gigiozzz/bundle-test-op"
		concat := "@"
		digest := "sha256:70ba938d4e11f219fc9dc0424e3e55173419a1da51598b341bb2162ea088a8a4"
	*/
	dir, err := ioutil.TempDir("/tmp", "crane-"+cr.Spec.Digest+"-")
	if err != nil {
		return nil, dir, err
	}

	err = bundles.ExtractImageTo(cr.Spec.Repository+"@"+cr.Spec.Digest, dir)
	if err != nil {
		return nil, dir, err
	}

	bundleDescriptor, err := bundles.ReadBundleDescriptor(dir + "/descriptor.yaml")
	if err != nil {
		return nil, dir, err
	}

	return bundleDescriptor.Components, dir, nil

}

func (bs *BundleService) retrieveSignatureImageRef(imageRef string) (string, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return "", err
	}

	ociremoteOpts := []ociremote.Option{}
	tag, err := ociremote.SignatureTag(ref, ociremoteOpts...)
	if err != nil {
		return "", err
	}

	return tag.Name(), nil
}

func (bs *BundleService) verifySignature(imageRef string, key string) error {
	/*
		o := &options.VerifyOptions{}

		annotations, err := o.AnnotationsMap()
		if err != nil {
			return err
		}

		hashAlgorithm, err := o.SignatureDigest.HashAlgorithm()
		if err != nil {
			return err
		}
	*/

	v := verify.VerifyCommand{
		KeyRef:         key,
		SkipTlogVerify: true,
	}

	err := v.Exec(context.TODO(), []string{imageRef})

	return err
}
