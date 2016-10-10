# ArduinoMQTT Snap

Snap that reads lines from the arduino and puts them into the MQTT topic sensor/arduino/in.

Build:
git clone ...
snapcraft

Run:

sudo snap install arduinomqtt_0.1_... --force-dangerous --devmode

This will publish to localhost:1833 to the MQTT topic sensor/arduino/in the data read that comes from the Arduino.
