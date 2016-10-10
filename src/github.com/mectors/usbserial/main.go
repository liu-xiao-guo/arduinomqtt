package main

import (
  "fmt"
  "log"
  "bufio"
  "flag"
	MQTT "github.com/eclipse/paho.mqtt.golang"
  "github.com/jacobsa/go-serial/serial"
)
var conns = flag.Int("conns", 10, "how many conns (0 means infinite)")
var host = flag.String("host", "localhost:1883", "hostname of broker")
var clientID = flag.String("clientid", "rfid", "the mqtt clientid")
var user = flag.String("user", "", "username")
var pass = flag.String("pass", "", "password")
var usb = flag.String("usb", "/dev/ttyUSB0", "usb device to read from")

var topic = "discovery"
var intopic = "sensor/arduino/in"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	flag.Parse()

  // Prepare the local MQTT connection
	opts2 := MQTT.NewClientOptions().AddBroker("tcp://"+*host)
  opts2.SetClientID(*clientID)


  //create and start a client using the above ClientOptions
  Clocal := MQTT.NewClient(opts2)
  if token := Clocal.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }
	fmt.Println("Connected locally")
  defer Clocal.Disconnect(250)

	// Say we are ready for action
	Clocal.Publish(topic, 0, false,	intopic)
  fmt.Println("Published to:"+topic+" that we are listening on:"+intopic)
    // Set up options.
    options := serial.OpenOptions{
      PortName: *usb,
      BaudRate: 9600,
      DataBits: 8,
      StopBits: 1,
      MinimumReadSize: 4,
    }
    // Open the port.
    port, err := serial.Open(options)
    if err != nil {
      log.Fatalf("serial.Open: %v", err)
    }

    // Make sure to close it later.
    defer port.Close()

    reader := bufio.NewReader(port)
    scanner := bufio.NewScanner(reader)

    scanner.Split(bufio.ScanLines)
    // first scan might have partial data so ignore
    scanner.Scan()
    for scanner.Scan() {
      if !Clocal.IsConnected() {
        if token := Clocal.Connect(); token.Wait() && token.Error() != nil {
          panic(token.Error())
        }
      }
      if ptoken := Clocal.Publish(intopic, 0, false,	scanner.Text()); ptoken.Wait() && ptoken.Error() != nil {
        panic(ptoken.Error())
      }
    }
}
