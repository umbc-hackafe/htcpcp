#include "teensy3/WProgram.h"
#include "teensy3/core_pins.h"
#include "teensy3/usb_serial.h"

#define BUFLEN 16

#define PIN_TURNTABLE_SERVO 1
#define PIN_KCUP_LOADER 2
#define PIN_KCUP_EJECTOR 3
#define PIN_TRAY_PISTON 4

const int ERR_UNKNOWN = -1;
const int ERR_NYI = -2;

const int ERR_KCUP_UNKNOWN = -100;
const int ERR_KCUP_EMPTY = -101;
const int ERR_KCUP_TRAY_FULL = -102;
const int ERR_KCUP_TRAY_EMPTY = -103;

const int ERR_COFFEE_DRY = -201;
const int ERR_KETTLE_DRY = -202;

const int ERR_ALREADY_BREWING = -301;

const int ERR_MUG_INVALID = -401;

enum Addin {
  SUGAR,
  MILK,
  HONEY,
}

char serial_buffer[BUFLEN] = {};
long serial_index = 0;

#define CMD_HEADER_A 0xCA
#define CMD_HEADER_B 0xFE

#define CMD_mug_seek 'a'
#define CMD_mug_get 'b'
#define CMD_kcup_tray_open 'c'
#define CMD_kcup_tray_close 'd'
#define CMD_kcup_load 'e'
#define CMD_kcup_eject 'f'
#define CMD_kcup_count 'g'
#define CMD_brew 'h'
#define CMD_addin_insert 'i'

int arglen[256] = {};
arglen[CMD_mug_seek] = 1;
arglen[CMD_mug_get] = 0;
arglen[CMD_kcup_tray_open] = 0;
arglen[CMD_kcup_tray_close] = 0;
arglen[CMD_kcup_load] = 0;
arglen[CMD_kcup_eject] = 0;
arglen[CMD_kcup_count] = 0;
arglen[CMD_brew] = 0;
arglen[CMD_addin_insert] = 2;

void send_data(int result) {

}

void handle_serial {
  if (Serial.available()) {
    int data = Serial.read();
    if (data == CMD_HEADER_A) {
      serial_index = 0;
    } else {
      serial_buffer[serial_index++] = data;
    }
  }

  if (serial_buffer[0] == CMD_HEADER_B && serial_index > 1) {
    if (serial_index > 1 + arglen[serial_buffer[1]]) {
      int res;
      switch (serial_buffer[1]) {
      case CMD_mug_seek:
	// int
	res = mug_seek(serial_buffer[2]);
	break;
      case CMD_mug_get:
	res = mug_get();
	break;
      case CMD_kcup_tray_open:
	res = kcup_tray_open();
	break;
      case CMD_kcup_tray_close:
	res = kcup_tray_close();
	break;
      case CMD_kcup_load:
	res = kcup_load();
	break;
      case CMD_kcup_eject:
	res = kcup_eject();
	break;
      case CMD_kcup_count:
	res = kcup_count();
	break;
      case CMD_brew:
	res = brew();
	break;
      case CMD_addin_insert:
	// int, int
	res = brew(serial_buffer[2], serial_buffer[3]);
	break;
      default:
	res = ERR_UNKNOWN;
	break;
      }
      send_data(res);
      serial_index = 0;
    }
  }
}

// Tell the turntable to seek to a particular mug
// Error if !(0 <= index <= 3)
int mug_seek(int index) {
  return ERR_NYI;
}

// Gets the current selected mug index, or the one that is currently being seeked
// Never returns an error
int mug_get() {
  return ERR_NYI;
}

// Ensure the kcup tray is open
// Repeated calls will do nothing
// Never returns an error, ideally
int kcup_tray_open() {
  return ERR_NYI;
}

// Ensure the kcup tray is closed
// Repeated calls will do nothing
// Never returns an error, ideally
int kcup_tray_close() {
  return ERR_NYI;
}

// Loads a kcup from the dispenser
// ERR_KCUP_TRAY_FULL if there is already a kcup loaded
// Opens the tray if it is closed
int kcup_load() {
  kcup_tray_open();
  return ERR_NYI;
}

// Ejects a kcup from the dispenser
// Opens the tray if it is closed
// Never returns an error
int kcup_eject() {
  kcup_tray_open();
  return ERR_NYI;
}

// Returns the number of kcups loaded
// If exact numbers are not available, returns 1 if kcups are
// available and 0 if no kcups are available
int kcup_count() {
  return ERR_NYI;
}

// Begins a brew cycle
// ERR_KCUP_TRAY_EMPTY if no kcup is loaded
// ERR_COFFEE_DRY if water reservoir is empty
// ERR_ALREADY_BREWING if a brew is already in progress
int brew() {
  return ERR_NYI;
}

// Dispenses an addin to the current mug
// Will probably not return errors
int addin_insert(enum Addin type, int amount) {
  return ERR_NYI;
}

int kettle_set_temp(int temp) {

}

int kettle_dispense(int amount) {

}

int main() {
  Serial.begin(9600);

  while (1) {
    handle_serial();
  }
}
