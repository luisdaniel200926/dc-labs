package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"math"
)

type Point struct {
	X, Y float64
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//generatePoints array
func generatePoints(s string) ([]Point, error) {

	points := []Point{}

	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	vals := strings.Split(s, ",")
	if len(vals) < 2 {
		return []Point{}, fmt.Errorf("Point [%v] was not well defined", s)
	}

	var x, y float64

	for idx, val := range vals {

		if idx%2 == 0 {
			x, _ = strconv.ParseFloat(val, 64)
		} else {
			y, _ = strconv.ParseFloat(val, 64)
			points = append(points, Point{x, y})
		}
	}
	return points, nil
}

// Distance between two Points
func getDistance(p, q Point) float64 {
	return math.Sqrt(math.Pow(q.X - p.X, 2) + math.Pow(q.Y - p.Y, 2))
}

// getArea gets the area inside from a given shape
func getArea(points []Point) float64 {
	//Calculating Area
	//https://www.mathopenref.com/coordpolygonarea.html

	val1 := 0.0
	val2 := 0.0
	for i := 0; i < len(points); i++ {
		if i == len(points)-1{
			val1 += points[len(points)-1].X * points[0].Y
			val2 += points[len(points)-1].Y * points[0].X
		}else{
			val1 += points[i].X * points[i+1].Y
			val2 += points[i].Y * points[i+1].X
		}
	}
	total :=math.Abs(val1 - val2) / 2
	return total

}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {

	sum := 0.0
	for i := 0; i < len(points); i++ {
		if i == len(points)-1{
			sum += getDistance(points[len(points)-1], points[0])
		}else{
		sum += getDistance(points[i], points[i+1])
		}
	}
	return sum

}

func onSegment(p, q, r Point) bool {
	if q.X <= math.Max(p.X, r.X) && q.X >= math.Min(p.X, r.X) &&
		q.Y <= math.Max(p.Y, r.Y) && q.Y >= math.Min(p.Y, r.Y) {
		return true
	}
	return false
}

func orientation(p, q, r Point) int {

	val :=0.0
	val = ((q.Y - p.Y) * (r.X - q.X)) - ((q.X - p.X) * (r.Y - q.Y))

    if (val > 0) {
        return 1 // Clockwise orientation
    } else if (val < 0) {
        return 2 // Counterclockwise orientation
    } else {
        return 0 // Colinear orientation
    }
}

func doIntersect(p1, q1, p2, q2 Point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)
	if o1 != o2 && o3 != o4 { //General Case
		return true
	} else if o1 == 0 && onSegment(p1, p2, q1) { // Special Cases
		return true
	} else if o2 == 0 && onSegment(p1, q2, q1) {// Special Cases
		return true
	} else if o3 == 0 && onSegment(p2, p1, q2) {// Special Cases
		return true
	} else if o4 == 0 && onSegment(p2, q1, q2) {// Special Cases
		return true
	}
	return false
}

func hasIntersections(points []Point) bool {
		for i := 0; i < len(points)-1; i++ {
			if i == len(points)-2{
				p1, q1 := points[i],points[i+1]
				p2, q2 := points[0], points[1]
				if doIntersect(p1, q1, p2, q2) {
					return true
				}
			}else{
				p1, q1 := points[i],points[i+1]
				for j := i + 2; j < len(points)-1; j++ {
					p2, q2 := points[j], points[j+1]
					if doIntersect(p1, q1, p2, q2) {
						return true
					}
				}
			}
		}
		return false
}


// handler handles the web request and reponds it
func handler(w http.ResponseWriter, r *http.Request) {

	var vertices []Point
	for k, v := range r.URL.Query() {
		if k == "vertices" {
			points, err := generatePoints(v[0])
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("error: %v", err))
				return
			}
			vertices = points
			break
		}
	}

	// Results gathering
	area := getArea(vertices)
	perimeter := getPerimeter(vertices)

	// Logging in the server side
	log.Printf("Received vertices array: %v", vertices)

	// Response construction
	response := fmt.Sprintf("Welcome to the Remote Shapes Analyzer\n")
	response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", len(vertices))

	if len(vertices)>2{
	response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
	response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
	response += fmt.Sprintf(" - Area            : %v\n", area)
	}else if hasIntersections(vertices) && len(vertices)>=4{
		response += fmt.Sprintf(" - Error -Your shape has some kind of intersections between lines-")

	}else{
		response += fmt.Sprintf(" - Error -Your shape has not enough vertex-")
	}
	// Send response to client
	fmt.Fprintf(w, response)
}
