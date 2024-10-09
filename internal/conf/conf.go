package conf

import (
	"context"
	"flag"
	"github.com/go-faker/faker/v4"
	"github.com/hashicorp/vault-client-go"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"github.com/mitchellh/mapstructure"
	"log"
	"os"
	"time"
)

func openVaultClient[T any](path, mount string, output T) error {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "http://localhost:8201"
	}
	vaultToken := os.Getenv("VAULT_SECRET")
	if vaultToken == "" {
		vaultToken = "secret"
	}

	client, err := vault.New(
		vault.WithAddress(vaultAddr),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return err
	}

	if err = client.SetToken(vaultToken); err != nil {
		return err
	}

	s, err := client.Secrets.KvV2Read(context.Background(), path, vault.WithMountPath(mount))
	if err != nil {
		log.Fatal(err)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  output,
		TagName: "json",
	})
	if err != nil {
		return err
	}
	err = decoder.Decode(s.Data.Data)
	if err != nil {
		return err
	}

	return nil
}

func LoadOtelConf() *secret_proto.Otel {
	otelConf := &secret_proto.Otel{}
	if flag.Lookup("test.v") != nil {
		err := faker.FakeData(&otelConf)
		util.Panic(err)
		return otelConf
	}

	err := openVaultClient("otel", "kv", otelConf)
	util.Panic(err)
	return otelConf
}

func LoadMinioConf() *secret_proto.Minio {
	minioConf := &secret_proto.Minio{}
	if flag.Lookup("test.v") != nil {
		err := faker.FakeData(&minioConf)
		util.Panic(err)
		return minioConf
	}
	err := openVaultClient("minio", "kv", minioConf)
	util.Panic(err)
	return minioConf
}

func LoadRabbitMQConf() *secret_proto.RabbitMQ {
	rabbitConf := &secret_proto.RabbitMQ{}
	if flag.Lookup("test.v") != nil {
		err := faker.FakeData(&rabbitConf)
		util.Panic(err)
		return rabbitConf
	}
	err := openVaultClient("rabbitmq", "kv", rabbitConf)
	util.Panic(err)
	return rabbitConf
}

func LoadJwtConf() *secret_proto.Jwt {
	jwtConf := &secret_proto.Jwt{}
	if flag.Lookup("test.v") != nil {
		err := faker.FakeData(&jwtConf)
		util.Panic(err)
		return jwtConf
	}
	err := openVaultClient("jwt", "kv", jwtConf)
	util.Panic(err)
	return jwtConf
}

func LoadAppConf() *AppConfig {
	appConf := &AppConfig{}
	if flag.Lookup("test.v") != nil {
		err := faker.FakeData(&appConf)
		util.Panic(err)
		return appConf
	}
	err := openVaultClient("user-service", "kv", appConf)
	util.Panic(err)
	return appConf
}
