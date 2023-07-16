package services

import (
	"fmt"
	"testing"
)

func TestRetrieveSignatureImageRef(t *testing.T) {
	bs := &BundleService{}

	data := "docker.io/gigiozzz/bundle-test-op@sha256:a41dbb9b16f052f1d26a22a5de34671e831cfb6fd327726f89bed5f8798dfd23"
	expectedName := "index.docker.io/gigiozzz/bundle-test-op:sha256-a41dbb9b16f052f1d26a22a5de34671e831cfb6fd327726f89bed5f8798dfd23.sig"
	actualName, err := bs.retrieveSignatureImageRef(data)
	if err != nil || actualName != expectedName {
		t.Fatalf("Invalid generation for %q. Expected %q, got %q error %s", data, expectedName, actualName, err)
	}
}

func TestVerifySignature(t *testing.T) {
	bs := &BundleService{}

	image := "docker.io/gigiozzz/bundle-test-op@sha256:a41dbb9b16f052f1d26a22a5de34671e831cfb6fd327726f89bed5f8798dfd23"
	k8sKeySecret := "k8s://bundle-operator/bundle-a4e2c0a3-key-secret"
	err := bs.verifySignature(image, k8sKeySecret)
	if err != nil {
		t.Fatalf("Invalid signature for %q. error %s", image, err)
	}
}

func TestVerifySignatureError(t *testing.T) {
	bs := &BundleService{}

	image := "docker.io/gigiozzz/bundle-test-op@sha256:a41dbb9b16f052f1d26a22a5de34671e831cfb6fd327726f89bed5f8798dfd23"
	k8sKeySecret := "k8s://bundle-operator/bundle-a4e2c0a3-key-secretssss"
	err := bs.verifySignature(image, k8sKeySecret)
	if err == nil {
		t.Fatalf("Invalid signature for %q. error %s", image, err)
	}
	fmt.Println("the error is: " + err.Error())
}
