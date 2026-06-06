package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
)

func TestAssetService_CreateAndGet(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "assetuser")

	expDate := "2027-01-01"
	resp, err := svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name: "example.com", AssetType: "domain", ExpireDate: &expDate,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Name != "example.com" {
		t.Errorf("name = %q, want %q", resp.Name, "example.com")
	}

	got, err := svc.AssetService.GetByID(userID, resp.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "example.com" {
		t.Errorf("got name = %q, want %q", got.Name, "example.com")
	}
}

func TestAssetService_Update(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "assetuser2")

	created, _ := svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name: "myserver", AssetType: "server",
	})

	updated, err := svc.AssetService.Update(userID, created.ID, &dto.UpdateAssetRequest{
		Name: "myserver-v2", AssetType: "server",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "myserver-v2" {
		t.Errorf("name = %q, want %q", updated.Name, "myserver-v2")
	}
}

func TestAssetService_Delete(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "assetuser3")

	created, _ := svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name: "to-delete", AssetType: "ssl_certificate",
	})

	err := svc.AssetService.Delete(userID, created.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = svc.AssetService.GetByID(userID, created.ID)
	if err == nil {
		t.Error("expected error getting deleted asset")
	}
}

func TestAssetService_UserIsolation(t *testing.T) {
	svc := newTestServices(t)
	userA := registerTestUser(t, svc, "userA_asset")
	userB := registerTestUser(t, svc, "userB_asset")

	created, _ := svc.AssetService.Create(userA, &dto.CreateAssetRequest{
		Name: "private-domain", AssetType: "domain",
	})

	_, err := svc.AssetService.GetByID(userB, created.ID)
	if err == nil {
		t.Error("user B should not access user A's asset")
	}

	err = svc.AssetService.Delete(userB, created.ID)
	if err == nil {
		t.Error("user B should not delete user A's asset")
	}
}
