package main

import (
	"testing"
)

func TestIsVilrosAvailable__yields_false_if_product_unavailable(t *testing.T) {
	// Given
	item := Item{
		Store: "vilros",
		Link:  "https://vilros.com/products/raspberry-pi-4-2gb-ram?src=raspberrypi",
		Ram:   2,
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
		Ram:   2,
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
