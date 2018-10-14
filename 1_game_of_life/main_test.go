package main

import (
	"reflect"
	"testing"
)

func TestGetSizes(t *testing.T) {
	for _, tc := range []struct {
		m    [][]int
		w, h int
	}{
		{nil, 0, 0},
		{[][]int{}, 0, 0},
		{[][]int{[]int{1}}, 1, 1},
	} {
		gotW, gotH := getSizes(tc.m)
		if gotW != tc.w || gotH != tc.h {
			t.Errorf(
				"getSizes(%v) = %d, %d, but want %d, %d",
				tc.m, gotW, gotH, tc.w, tc.h,
			)
		}
	}
}

func TestPlaceFigure(t *testing.T) {
	validArgs := struct {
		dest, fig [][]int
		x, y      int
	}{
		dest: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		fig: [][]int{
			{1, 2},
			{3, 4},
		},
		x: 1,
		y: 1,
	}
	want := [][]int{
		{0, 0, 0},
		{0, 1, 2},
		{0, 3, 4},
	}

	got, err := placeFigure(validArgs.dest, validArgs.fig, validArgs.x, validArgs.y)
	if err != nil {
		t.Error(err)
	}

	if ok := reflect.DeepEqual(got, want); !ok {
		t.Errorf("placeFigure() = %v, want %v", got, want)
	}
}

func TestCountNeighbours(t *testing.T) {
	m := [][]int{
		{1, 0, 0},
		{1, 1, 0},
		{0, 0, 1},
	}

	got := countNeighbours(m, 1, 1)
	want := 3
	if got != want {
		t.Errorf("countNeighbours() = %d, expected %d", got, want)
	}
}
