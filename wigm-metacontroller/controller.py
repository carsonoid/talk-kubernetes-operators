from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import json

class Controller(BaseHTTPRequestHandler):
  def sync(self, parent, children):
    name = parent["metadata"]["name"]

    # default status
    status = {
      "deployment": {
        "created": bool(children["Deployment.apps/v1"]),
      },
      "service": {
        "created": bool(children["Service.v1"]),
        "type": children["Service.v1"].get("wigm-host-%s"%(name), {}).get("spec", {}).get("type","")
      },
      "ingress": {
        "created": bool(children["Ingress.extensions/v1beta1"]),
      }
    }

    # Generate the desired child objects
    children = [
      self.get_deployment(parent, name),
      self.get_service(parent, name),
    ]

    #conditionally create ingress, default TRUE
    if parent["spec"].get("ingress", {}).get("enabled", True):
      children.append(
        self.get_ingress(parent, name),
      )

    return {"status": status, "children": children}

  def do_POST(self):
    # Serve the sync() function as a JSON webhook.
    observed = json.loads(self.rfile.read(int(self.headers.getheader("content-length"))))
    desired = self.sync(observed["parent"], observed["children"])

    self.send_response(200)
    self.send_header("Content-type", "application/json")
    self.end_headers()
    self.wfile.write(json.dumps(desired))

  def get_deployment(self, parent, name):
    gifname = parent["spec"]["gif"].get("name", name)
    giflink = parent["spec"]["gif"]["link"]

    deployment = {
      "apiVersion": "apps/v1",
      "kind": "Deployment",
      "metadata": {
        "name": "wigm-host-%s"%(name),
        "labels": {
          "app": "wigm",
          "gif": name
        }
      },
      "spec": {
        "replicas": 2,
        "selector": {
          "matchLabels": {
            "app": "wigm",
            "gif": name
          }
        },
        "template": {
          "metadata": {
            "labels": {
              "app": "wigm",
              "gif": name
            }
          },
          "spec": {
            "containers": [
              {
                "name": "gifhost",
                "image": "nginx:1.15-alpine",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ],
                "env": [
                  {
                    "name": "GIF_NAME",
                    "value": gifname
                  },
                  {
                    "name": "GIF_SOURCE_LINK",
                    "value": giflink
                  }
                ],
                "command": [
                  "sh",
                  "-exc",
                  START_SCRIPT
                ]
              }
            ]
          }
        }
      }
    }

    return deployment

  def get_service(self, parent, name):
    service = {
      "kind": "Service",
      "apiVersion": "v1",
      "metadata": {
        "name": "wigm-host-%s"%(name),
        "labels": {
          "app": "wigm",
          "gif": name
        }
      },
      "spec": {
        "selector": {
          "app": "wigm",
          "gif": name
        },
        "ports": [
          {
            "protocol": "TCP",
            "port": 80,
            "targetPort": 80
          }
        ]
      }
    }

    # conditionally set service kind
    if parent["spec"].get("service", {}).get("create_cloud_lb"):
      service["spec"]["type"] = "LoadBalancer"

    return service


  def get_ingress(self, parent, name):
    ingress = {
      "apiVersion": "extensions/v1beta1",
      "kind": "Ingress",
      "metadata": {
        "name": "wigm-host-%s"%(name),
        "labels": {
          "app": "wigm",
          "gif": name
        }
      },
      "spec": {
        "rules": [
          {
            "host": "%s.wigm.carson-anderson.com"%(name),
            "http": {
              "paths": [
                {
                  "path": "/",
                  "backend": {
                    "serviceName": "wigm-host-%s"%(name),
                    "servicePort": 80
                  }
                }
              ]
            }
          }
        ]
      }
    }

    return ingress

START_SCRIPT = '''
# get curl
apk update && apk add curl

# go to data dir
cd /usr/share/nginx/html

# fetch original
curl -Lo wigm.gif "$GIF_SOURCE_LINK"

# write page using heredoc
cat > index.html <<EOF
<head>
  <title>WIGM: $GIF_NAME</title>
</head>
<body>
  <h1>WIGM: $GIF_NAME</h1>
  <img src="wigm.gif" />
</body>
EOF

# start nginx
exec nginx -g "daemon off;"
'''

HTTPServer(("", 80), Controller).serve_forever()
