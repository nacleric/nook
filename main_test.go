package main

import (
	"testing"
)

func TestIsVilrosAvailable__yields_false_if_product_unavailable(t *testing.T) {
	// Given
	item := Item{
		Store: "vilros",
		Link:  "https://vilros.com/products/raspberry-pi-4-2gb-ram?src=raspberrypi",
		Ram:   42,
		Misc:  "https://vilros.com/products/raspberry-pi-4-2gb-ram.js",
	}
	expected_res := false

	// When
	sut, _ := isVilrosAvailable(item)
	res := sut

	// Then
	if res != expected_res {
		t.Errorf("left == %v, right == %v", res, expected_res)
	}
}

func TestIsAdaFruitAvailable__yields_false_if_product_unavailable(t *testing.T) {
	// Given
	item := Item{
		Store: "adafruit",
		Link:  "https://www.adafruit.com/product/4292",
		Ram:   42,
	}
	expected_res := false

	// When
	sut, _ := isAdaFruitAvailable(item)
	res := sut

	// Then
	if res != expected_res {
		t.Errorf("left == %v, right == %v", res, expected_res)
	}
}

func TestIsPiShopAvailable__yields_false_if_product_unavailable(t *testing.T) {
	// Given
	item := Item{
		Store: "pishop.us",
		Link:  "https://www.pishop.us/product/raspberry-pi-4-model-b-2gb/?src=raspberrypi",
		Ram:   42,
	}
	expected_res := false

	// When
	sut, _ := isPiShopAvailable(item)
	res := sut

	// Then
	if res != expected_res {
		t.Errorf("left == %v, right == %v", res, expected_res)
	}

}

func TestNotifyEric(t *testing.T) {
	// Given
	dg := initBot()

	// When
	// Then
	notifyEric(dg, "Test NotifyEric() passed")
	dg.Close()
}
