meta {
    version = "v1.0.0"
    description = "A central place for all your Portmaster plugins"
}

plugin "yaegi" {
        source = "https://github.com/ppacher/portmaster-plugin-yaegi"
        version = "v0.0.1"
        pluginTypes = [
            "decider"
        ]
        author = "Patrick Pacher"
        license = "GPL"
        description = "Use the Go language to write custom rules"
        tags = ["rule-as-code"]

        artifact "linux" {
            amd64 = ""
        }
}

plugin "prometheus" {
        source = "https://github.com/ppacher/portmaster-plugin-prometheus"
        version = "v0.0.1"
        pluginTypes = [
            "reporter"
        ]
        author = "Patrick Pacher"
        license = "GPL"
        description = "Export connection metrics to prometheus"
        artifact "linux" {
            amd64 = ""
        }
}

plugin "hostsfile" {
        source = "https://github.com/ppacher/portmaster-plugin-hosts"
        version = "v0.0.1"
        pluginTypes = [
            "resolver"
        ]
        author = "Patrick Pacher"
        license = "GPL"
        description = "Add support for /etc/hosts"
        artifact "linux" {
            amd64 = ""
        }
}

plugin "dnscrypt-client" {
        source = "https://github.com/ppacher/portmaster-plugin-dnscrypt"
        version = "v0.0.1"
        pluginTypes = [
            "resolver"
        ]
        author = "Patrick Pacher"
        license = "GPL"
        description = "Add support to use DNSCrypt servers"
        artifact "linux" {
            amd64 = ""
        }
}