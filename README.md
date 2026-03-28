# translate-proxy

轻量级翻译代理，把翻译 API（如阿里云 qwen-mt-flash）包装成标准 OpenAI 兼容接口，可直接接入 [Open WebUI](https://github.com/open-webui/open-webui) 使用。

Go 实现，运行内存约 5MB。

## 快速启动

```bash
git clone https://github.com/ladydd/translate-proxy.git
cd translate-proxy
```

编辑 `config.json`，填入你的 API Key：

```json
{
  "port": "8787",
  "models": {
    "qwen-mt-flash": {
      "api_key": "sk-你的key",
      "api_base": "https://dashscope.aliyuncs.com/compatible-mode/v1",
      "model": "qwen-mt-flash",
      "timeout": 30,
      "translation_options": {
        "source_lang": "auto",
        "target_lang": "Chinese"
      }
    }
  }
}
```

启动：

```bash
docker compose up -d --build
```

## 添加模型

在 `config.json` 的 `models` 里加一条就行：

```json
{
  "port": "8787",
  "models": {
    "qwen-mt-flash": {
      "api_key": "sk-xxx",
      "api_base": "https://dashscope.aliyuncs.com/compatible-mode/v1",
      "model": "qwen-mt-flash",
      "timeout": 30,
      "translation_options": { "source_lang": "auto", "target_lang": "Chinese" }
    },
    "qwen-mt-plus": {
      "api_key": "sk-yyy",
      "api_base": "https://dashscope.aliyuncs.com/compatible-mode/v1",
      "model": "qwen-mt-plus",
      "timeout": 30,
      "translation_options": { "source_lang": "auto", "target_lang": "Chinese" }
    }
  }
}
```

重新构建即可：`docker compose up -d --build`

## 接入 Open WebUI

1. 管理员设置 → 外部连接
2. 添加 OpenAI API 连接：
   - API Base URL：`http://translate-proxy:8787/v1`（Docker 容器间）或 `http://localhost:8787/v1`
   - API Key：随便填
3. 模型列表里会出现 `config.json` 中配置的所有模型

## Docker 网络

Open WebUI 也是 Docker 的话，两个容器需要在同一网络，在 `docker-compose.yml` 里加：

```yaml
services:
  translate-proxy:
    build: .
    container_name: translate-proxy
    ports:
      - "8787:8787"
    restart: unless-stopped
    networks:
      - 你的网络名

networks:
  你的网络名:
    external: true
```



## 效果展示

![接入配置](https://img.weini.xin/d/photos/obsidian/20260328112058264.png?sign=JC8YiwpAs3T__DN7BTYT_txXG6YvsgifmgegTEurLD0=:0)

![翻译效果](https://img.weini.xin/d/photos/obsidian/20260328112136184.png?sign=ib4cRIPzdIs7pmr1aSfJo1z2CC0nNSYUfMx0bQf5A7U=:0)

## 协议

MIT
