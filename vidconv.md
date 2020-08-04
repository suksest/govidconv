---
title: Govidconv v1.0
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: false
highlight_theme: darkula
headingLevel: 2
generator: widdershins v4.0.1

---

<h1 id="govidconv">Govidconv v1.0</h1>

Base URLs:

* <a href="http://localhost:5000">http://localhost:5000</a>

Email: <a href="mailto:sukmasetyaji@gmail.com">Support</a> 

<h1 id="govidconv-video">Video</h1>


`POST /video/convert`

> Body parameter

```yaml
file: /home/user/Videos/lovely_cat.mp4
compression: 1
format: string
video_codec: string
audio_codec: string

```

<h3 id="convert-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|true|none|
|» file|body|string(binary)|true|none|
|» compression|body|number(float)|true|none|
|» format|body|string(extension)|true|none|
|» video_codec|body|string(video_codec)|true|none|
|» audio_codec|body|string(audio_codec)|true|none|

> Example responses

> 200 Response

```json
{
  "downloadPath": "http://localhost:5000/tmp/1UjfE___a014b86ac8b77e433090d45edd89e6c83ae94fde46c5ade80f19dd9f4f2c540d.mp4"
}
```

<h3 id="convert-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|video converted successfully|[ConvertResponse](#schemaconvertresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|video fail to convert due to empty request body or bad input|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|server error due to fail in operating system operations|None|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_ConvertResponse">ConvertResponse</h2>

<a id="schemaconvertresponse"></a>
<a id="schema_ConvertResponse"></a>
<a id="tocSconvertresponse"></a>
<a id="tocsconvertresponse"></a>

```json
{
  "downloadPath": "http://localhost:5000/tmp/1UjfE___a014b86ac8b77e433090d45edd89e6c83ae94fde46c5ade80f19dd9f4f2c540d.mp4"
}

```

Convert Response

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|downloadPath|string|true|none|none|

