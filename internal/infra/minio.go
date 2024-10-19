package infra

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/secret_proto"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func NewMinio(cred *secret_proto.Minio) *minio.Client {
	minioClient, err := minio.New(cred.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cred.AccessId, cred.SecretAccessKey, ""),
		Secure: cred.UseSsl,
	})
	collection.PanicIfErr(err)

	exist, err := minioClient.BucketExists(context.Background(), cred.PrivateBucket)
	collection.PanicIfErr(err)

	if !exist {
		err = minioClient.MakeBucket(context.Background(), cred.PrivateBucket, minio.MakeBucketOptions{})
		collection.PanicIfErr(err)
	}

	log.Info().Msg("initialization minio successfully")
	return minioClient
}
