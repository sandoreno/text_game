package main

import (
	"strings"
)

type Player struct {
	location string   // текущая локация игрока
	backpack bool     // наличие рюкзака
	items    []string // наличие предметов (нет доступа без рюкзака)
}

type Room struct {
	name               string          // название комнаты
	next_room          []string        // в какие комнаты можно попасть
	items              map[string]bool // какие предметы лежат в комнате
	access             bool            //доступность комнаты
	title              string
	activity           map[string]string
	check_target       func(item, target string) string
	create_description func(room Room) string
}

// создадим игрока и игровой мир
var player *Player
var world map[string]Room

func initGame() {
	world = map[string]Room{
		"кухня": {
			name:      "кухня",
			next_room: []string{"коридор"},
			items:     map[string]bool{"чай": true},
			access:    true,
			title:     "кухня, ничего интересного",
			create_description: func(room Room) string {
				return "ты находишься на кухне, " + descriptionOfItem(room.items) + generateDescriptionKitchen() + descriptionOfPaths(room.next_room)
			},
		},
		"коридор": {
			name:      "коридор",
			next_room: []string{"кухня", "комната", "улица"},
			items:     map[string]bool{},
			access:    true,
			title:     "ничего интересного",
			activity:  map[string]string{"дверь": "ключи"},
			check_target: func(item, target string) string {
				if world[player.location].activity[target] != item {
					return "не к чему применить"
				}
				street := world["улица"]
				street.access = true
				world["улица"] = street
				return "дверь открыта"
			},
			create_description: func(room Room) string {
				return "ничего интересного" + descriptionOfItem(room.items) + descriptionOfPaths(room.next_room)
			},
		},
		"комната": {
			name:      "комната",
			next_room: []string{"коридор"},
			items:     map[string]bool{"ключи": true, "конспекты": true, "рюкзак": true},
			access:    true,
			title:     "ты в своей комнате",
			create_description: func(room Room) string {
				return descriptionOfItem(room.items) + descriptionOfPaths(room.next_room)
			},
		},
		"улица": {
			name:      "улица",
			next_room: []string{"домой"},
			items:     map[string]bool{"чай": true},
			access:    false,
			title:     "на улице весна",
			create_description: func(room Room) string {
				return "на улице весна." + descriptionOfItem(room.items) + descriptionOfPaths(room.next_room)
			},
		},
	}
	player = &Player{
		location: "кухня",
		backpack: false,
		items:    []string{},
	}
}

func generateDescriptionKitchen() string {
	if player.backpack {
		return ", надо идти в универ"
	}
	return ", надо собрать рюкзак и идти в универ"
}

func descriptionOfItem(items map[string]bool) string {
	var table []string
	chair := ""

	for _, item := range []string{"ключи", "конспекты", "чай"} {
		if items[item] {
			table = append(table, item)
		}
	}

	if items["рюкзак"] {
		chair = "на стуле: рюкзак"
	}

	var parts []string
	if len(table) > 0 {
		parts = append(parts, "на столе: "+strings.Join(table, ", "))
	}
	if chair != "" {
		parts = append(parts, chair)
	}

	if len(parts) == 0 {
		return "пустая комната"
	}
	return strings.Join(parts, ", ")
}

func descriptionOfPaths(rooms []string) string {
	return ". можно пройти - " + strings.Join(rooms, ", ")
}

func transition(room string) string {
	// получаем текущую комнату и проверяем на наличие пути
	current_room := world[player.location]
	// проверяем доступность комнаты
	if !world[room].access {
		return "дверь закрыта"
	}

	description := ""
	// флаг наличия перехода
	is_valid := false
	for _, value := range current_room.next_room {
		if value == room {
			is_valid = true
			break
		}
	}
	// если нет возможности пройти
	if !is_valid {
		description += "нет пути в " + room
	} else {
		// меняем комнату игрока
		player.location = room
		description += world[player.location].title + descriptionOfPaths(world[player.location].next_room)
	}

	return description
}

func apply(item, target string) string {

	if !player.backpack {
		return "нет предмета в инвентаре - " + item
	}

	isInventory := false
	for _, value := range player.items {
		if item == value {
			isInventory = true
			break
		}
	}
	if !isInventory {
		return "нет предмета в инвентаре - " + item
	}

	return world[player.location].check_target(item, target)
}

func takeItem(item string) string {
	room := world[player.location]

	if !player.backpack {
		return "некуда класть"
	}
	if !room.items[item] {
		return "нет такого"
	}
	delete(room.items, item)
	player.items = append(player.items, item)

	return "предмет добавлен в инвентарь: " + item
}

// метод получения рюкзака
func getBackpack(item string) string {
	rooms := world[player.location]

	if !rooms.items[item] && item != "рюкзак" {
		return "нет рюкзака"
	}

	player.backpack = true
	delete(rooms.items, "рюкзак")
	return "вы надели: рюкзак"
}

func handleCommand(command string) string {
	com := strings.Fields(command)
	switch com[0] {
	case "осмотреться":
		return world[player.location].create_description(world[player.location])
	case "идти":
		return transition(com[1])
	case "применить":
		return apply(com[1], com[2])
	case "взять":
		return takeItem(com[1])
	case "надеть":
		return getBackpack(com[1])
	default:
		return "неизвестная команда"
	}
}
