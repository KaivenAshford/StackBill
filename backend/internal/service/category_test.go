package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
)

func TestCategoryService_CreateAndGet(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "catuser")

	resp, err := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "Custom Category", Type: "subscription", Color: "#ff0000", Icon: "star",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Name != "Custom Category" {
		t.Errorf("name = %q, want %q", resp.Name, "Custom Category")
	}

	got, err := svc.CategoryService.GetByID(userID, resp.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "Custom Category" {
		t.Errorf("got name = %q, want %q", got.Name, "Custom Category")
	}
}

func TestCategoryService_DuplicateName(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "catuser2")

	_, _ = svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "SameName", Type: "subscription",
	})

	_, err := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "SameName", Type: "subscription",
	})
	if err == nil {
		t.Fatal("expected error for duplicate category name")
	}
}

func TestCategoryService_Update(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "catuser3")

	created, _ := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "Original", Type: "asset", Color: "#000000",
	})

	updated, err := svc.CategoryService.Update(userID, created.ID, &dto.UpdateCategoryRequest{
		Name: "Updated", Type: "asset", Color: "#ffffff",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("name = %q, want %q", updated.Name, "Updated")
	}
}

func TestCategoryService_Delete(t *testing.T) {
	svc := newTestServices(t)
	userID := registerTestUser(t, svc, "catuser4")

	created, _ := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "ToDelete", Type: "subscription",
	})

	err := svc.CategoryService.Delete(userID, created.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = svc.CategoryService.GetByID(userID, created.ID)
	if err == nil {
		t.Error("expected error getting deleted category")
	}
}

func TestCategoryService_UserIsolation(t *testing.T) {
	svc := newTestServices(t)
	userA := registerTestUser(t, svc, "catuserA")
	userB := registerTestUser(t, svc, "catuserB")

	created, _ := svc.CategoryService.Create(userA, &dto.CreateCategoryRequest{
		Name: "PrivateCat", Type: "subscription",
	})

	_, err := svc.CategoryService.GetByID(userB, created.ID)
	if err == nil {
		t.Error("user B should not access user A's category")
	}
}
