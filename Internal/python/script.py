from mitmproxy import http
from kafka import KafkaProducer
import json

# Kafka Producer 配置
producer = KafkaProducer(
    bootstrap_servers='117.50.85.130:9092',  # 修改为你的 Kafka 地址
    value_serializer=lambda v: json.dumps(v).encode('utf-8')
)

def request(flow: http.HTTPFlow):
    """拦截 HTTP 请求并发送到 Kafka"""
    data = {
        "url": flow.request.url,
        "method": flow.request.method,
        "headers": dict(flow.request.headers),
        "content": flow.request.text
    }
    producer.send("mitmproxy-topic", {"type": "request", "data": data})
    print(f"📤 已发送请求到 Kafka: {data}")

def response(flow: http.HTTPFlow):
    """拦截 HTTP 响应并发送到 Kafka"""
    content_type = flow.response.headers.get("Content-Type", "").lower()

    # 只处理文本/JSON数据，避免发送图片等非文本数据
    if "text" in content_type or "json" in content_type:
        data = {
            "url": flow.request.url,
            "method": flow.request.method,
            "status_code": flow.response.status_code,
            "headers": dict(flow.response.headers),
            "content": flow.response.text
        }
        producer.send("mitmproxy-topic", {"type": "response", "data": data})
        print(f"✅ 已发送响应到 Kafka: {data}")
    else:
        print(f"⏩ 跳过非文本内容: {flow.request.url}")


# command：mitmdump --set stream_websocket=true -s script.py