package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Room represents a hotel room
type Room struct {
	Number     int
	Type       string
	Price      float64
	IsReserved bool
}

// reservation represents a booking
type Reservation struct {
	ID         int
	RoomNumber int
	GuestName  string
	GuestEmail string
	CheckIn    time.Time
	CheckOut   time.Time
	CreatedAt  time.Time
}

// HotelSystem manages the hotel operations
type HotelSystem struct {
	rooms        []Room
	reservations []Reservation
	nextResID    int
}

func NewHotelSystem() *HotelSystem {
	return &HotelSystem{
		rooms: []Room{
			{Number: 101, Type: "Single", Price: 80.00, IsReserved: false},
			{Number: 102, Type: "Double", Price: 100.00, IsReserved: false},
			{Number: 201, Type: "Suite", Price: 280.00, IsReserved: false},
			{Number: 103, Type: "Single", Price: 90.00, IsReserved: false},
			{Number: 202, Type: "Double", Price: 130.00, IsReserved: false},
			{Number: 301, Type: "Suite", Price: 289.00, IsReserved: false},
		},

		reservations: []Reservation{},
		nextResID:    1,
	}
}

func main() {
	hotel := NewHotelSystem()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("======================================================")
	fmt.Println("      Welcome to the new Hotel Revervation System")
	fmt.Println("======================================================")

	for {
		displayMenu()
		option := readInput(reader, "Select an option: ")

		switch option {
		case "1":
			viewAvailableRooms()
		case "2":
			makeReservation(reader)
		case "3":
			viewReservations()
		case "4":
			cancelReservations()
		case "5":
			searchRoom()
		case "6":
			fmt.Println("\nThank you for using our system. Goodbye...")
			return
		default:
			fmt.Println("Invalid option selected. Please select 1-6.")
		}
		fmt.Println()
	}
}

func displayMenu() {
	fmt.Println("--------------------------------------")
	fmt.Println("|      MAIN MENU                     |")
	fmt.Println("--------------------------------------")
	fmt.Println("1. View Available Rooms")
	fmt.Println("2. Make Reservation")
	fmt.Println("3. View Reservations")
	fmt.Println("4. Cancel Reservations")
	fmt.Println("5. Search Room by Type")
	fmt.Println("6. Exit")
	fmt.Println("--------------------------------------- ")

}

func viewAvailableRooms() {
	fmt.Println("Available Rooms:")
	for _, room := range rooms {
		if !room.IsReserved {
			fmt.Printf("Room Number: %d, Type: %s\n", room.Number, room.Type)
		}
	}
}

func makeReservation(reader *bufio.Reader) {
	viewAvailableRooms()
	fmt.Println("***************************")

	fmt.Print("Enter Room Number to reserve: ")
	roomInput, _ := reader.ReadString('\n')
	roomInput = strings.TrimSpace(roomInput)
	roomNumber, err := strconv.Atoi(roomInput)
	if err != nil {
		fmt.Println("Invalid room number.")
		return
	}

	roomIndex := -1
	for i, room := range rooms {
		if room.Number == roomNumber && !room.IsReserved {
			roomIndex = i
			break
		}
	}
	if roomIndex == -1 {
		fmt.Println("Room not available.")
		return
	}

	for _, room := range rooms {
		if room.Number == roomNumber {
			if room.IsReserved {
				fmt.Println("Room is already reserved.")
			} else {
				fmt.Print("Enter Guest Name: ")
				guestName, _ := reader.ReadString('\n')
				guestName = strings.TrimSpace(guestName)

				room.IsReserved = true
				reservations = append(reservations, Reservation{RoomNumber: roomNumber, GuestName: guestName})
				fmt.Printf("Room %d reserved successfully for %s.\n", roomNumber, guestName)
			}
			return
		}
	}
	fmt.Println("Room not found.")
}

func viewReservations() {
	if len(reservations) == 0 {
		fmt.Println("No reservations found.")
		return
	}

	fmt.Println("Current Reservations:")
	for _, res := range reservations {
		fmt.Printf("Room Number: %d, Guest Name: %s\n", res.RoomNumber, res.GuestName)
	}

}
