package internal

import (
	ics "github.com/arran4/golang-ical"
	"io"
	"os"
	"testing"
)

func getTestData(t *testing.T, name string) (string, *App) {
	f, err := os.Open("testdata/" + name)
	if err != nil {
		t.Fatal("can't open testdata")
	}
	all, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("can't read testdata")
	}
	app, err := newApp()
	if err != nil {
		t.Fatal("can't create test subject", err)
	}
	return string(all), app
}

func TestReplacement(t *testing.T) {
	r1 := Replacement{"b", "b"}
	if r1.isLessThan(&r1) {
		t.Error("Replacement should not be less than itself")
		return
	}
	if r1.isLessThan(&Replacement{key: "longer key first"}) {
		t.Error("Replacement should sort longer prefix first")
		return
	}
	if !r1.isLessThan(&Replacement{key: ""}) {
		t.Error("Replacement should sort longer prefix first")
		return
	}
	if r1.isLessThan(&Replacement{key: "a"}) {
		t.Error("Replacement should sort alphabetically")
		return
	}
	if !r1.isLessThan(&Replacement{key: "c"}) {
		t.Error("Replacement should sort alphabetically")
		return
	}
	if r1.isLessThan(&Replacement{key: r1.key, value: "a"}) {
		t.Error("Replacement with equal key should sort alphabetically by value")
		return
	}
	if !r1.isLessThan(&Replacement{key: r1.key, value: "c"}) {
		t.Error("Replacement with equal key should sort alphabetically by value")
		return
	}

}

func TestDeduplication(t *testing.T) {
	testData, app := getTestData(t, "duplication.ics")
	calendar, err := app.getCleanedCalendar([]byte(testData))
	if err != nil {
		t.Error(err)
		return
	}
	if len(calendar.Components) != 1 {
		t.Errorf("Calendar should have only 1 entry after deduplication but has %d", len(calendar.Components))
		return
	}
}

func TestNameShortening(t *testing.T) {
	testData, app := getTestData(t, "nameshortening.ics")
	calendar, err := app.getCleanedCalendar([]byte(testData))
	if err != nil {
		t.Error(err)
		return
	}
	summary := calendar.Components[0].(*ics.VEvent).GetProperty(ics.ComponentPropertySummary).Value
	if summary != "ERA" {
		t.Errorf("Einführung in die Rechnerarchitektur (IN0004) VO, Standardgruppe should be shortened to ERA but is %s", summary)
		return
	}
}

func TestLocationReplacement(t *testing.T) {
	testData, app := getTestData(t, "location.ics")
	calendar, err := app.getCleanedCalendar([]byte(testData))
	if err != nil {
		t.Error(err)
		return
	}
	location := calendar.Components[0].(*ics.VEvent).GetProperty(ics.ComponentPropertyLocation).Value
	expectedLocation := "Boltzmannstr. 15\\, 85748 Garching b. München"
	if location != expectedLocation {
		t.Errorf("Location should be shortened to %s but is %s", expectedLocation, location)
		return
	}
	desc := calendar.Components[0].(*ics.VEvent).GetProperty(ics.ComponentPropertyDescription).Value
	expectedDescription := "https://nav.tum.de/room/5508.02.801\\nMW 1801\\, Ernst-Schmidt-Hörsaal (5508.02.801)\\nEinführung in die Rechnerarchitektur\\nfix\\; Abhaltung\\;"
	if desc != expectedDescription {
		t.Errorf("Description should be %s but is %s", expectedDescription, desc)
		return
	}
}
