import random
import time
import json
import ssl
from paho.mqtt import client as mqtt


device_id = ""
sas_token = ""
iot_hub_name = "gps-trackly"

# Corrected: Using actual Azure IoT Hub hostname
hostname = f"{iot_hub_name}.azure-devices.net"


def on_connect(client, userdata, flags, rc):
    print("Device connected with result code: " + str(rc))
    if rc == 0:
        print("Connected to Azure IoT Hub successfully!")
    else:
        print(f"Connection failed with code {rc}")


def on_disconnect(client, userdata, rc):
    print("Device disconnected with result code: " + str(rc))


def on_publish(client, userdata, mid):
    print("Device sent message")


def on_subscribe(client, userdata, mid, granted_qos):
    print("Topic subscribed!")


def on_message(client, userdata, msg):
    print("Received message!")
    print("Topic: '" + msg.topic + "', payload: " + str(msg.payload))


def simulate_device():
    # Corrected: Using proper client ID format for Azure IoT Hub
    client = mqtt.Client(client_id=device_id, protocol=mqtt.MQTTv311, clean_session=False)
    
    client.on_connect = on_connect
    client.on_disconnect = on_disconnect
    client.on_publish = on_publish
    client.on_subscribe = on_subscribe
    client.on_message = on_message

    # Azure IoT Hub connection details - ENABLE THESE FOR ACTUAL AZURE CONNECTION
    print(f"Setting up connection to Azure IoT Hub: {hostname}")
    username = f"{hostname}/{device_id}/?api-version=2021-04-12"
    client.username_pw_set(username=username, password=sas_token)

    # TLS configuration for Azure IoT Hub
    client.tls_set(ca_certs=None, certfile=None, keyfile=None,
                   cert_reqs=ssl.CERT_REQUIRED, tls_version=ssl.PROTOCOL_TLSv1_2, ciphers=None)
    client.tls_insecure_set(False)

    # Corrected: Connecting to Azure IoT Hub instead of test broker
    print("Connecting to Azure IoT Hub...")
    # For actual Azure connection, use port 8883 (TLS)
    client.connect(hostname, port=8883)
    
    # Start the network loop
    client.loop_start()

    # Topic for sending device-to-cloud messages
    topic = f"devices/{device_id}/messages/events/"
    print(f"MQTT topic: {topic}")

    try:
        while True:
            # Generate GPS-like sensor data (similar to our GPS simulator)
            latitude = random.uniform(40.4774, 40.9176)   # NYC latitude range
            longitude = random.uniform(-74.2591, -73.7002)  # NYC longitude range
            
            # Send data to IoT hub
            send_data_to_iot_hub(client, latitude, longitude, topic)

            # Wait for some time before sending the next data
            time.sleep(10)
            
    except KeyboardInterrupt:
        print("Simulation stopped by user")
    finally:
        client.loop_stop()
        client.disconnect()


def send_data_to_iot_hub(device_client, latitude, longitude, topic):
    payload = {
        "device_id": device_id,
        "latitude": latitude,
        "longitude": longitude,
        "timestamp": time.time()
    }
    message = json.dumps(payload)
    result = device_client.publish(topic, message, qos=1)
    if result.rc == mqtt.MQTT_ERR_SUCCESS:
        print(f"Message sent: {message}")
    else:
        print(f"Failed to send message: {result.rc}")


# Start simulating the device
if __name__ == "__main__":
    simulate_device()