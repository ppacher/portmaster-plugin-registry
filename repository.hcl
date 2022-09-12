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
        license = "MIT"
        description = "Use the Go language to write custom rules"
        tags = ["rule-as-code"]

        artifact_template = "{{source}}/releases/download/{{version}}/{{archive_file}}_{{stripped_version}}_{{os}}_{{arch}}.tar.gz"
        archive_file = "portmaster-plugin-yaegi"
}

plugin "prometheus" {
        source = "https://github.com/ppacher/portmaster-plugin-prometheus"
        version = "v0.0.1"
        pluginTypes = [
            "reporter"
        ]
        author = "Patrick Pacher"
        license = "MIT"
        description = "Export connection metrics to prometheus"
        artifact_template = "{{source}}/releases/download/{{version}}/{{archive_file}}_{{stripped_version}}_{{os}}_{{arch}}.tar.gz"
        archive_file = "portmaster-plugin-prometheus"
}

plugin "hostsfile" {
        source = "https://github.com/ppacher/portmaster-plugin-hosts"
        version = "v0.0.1"
        pluginTypes = [
            "resolver"
        ]
        author = "Patrick Pacher"
        license = "MIT"
        description = "Add support for /etc/hosts"
        artifact_template = "{{source}}/releases/download/{{version}}/{{archive_file}}_{{stripped_version}}_{{os}}_{{arch}}.tar.gz"
        archive_file = "portmaster-plugin-hosts"
}

plugin "dnscrypt-client" {
        source = "https://github.com/ppacher/portmaster-plugin-dnscrypt"
        version = "v0.0.1"
        pluginTypes = [
            "resolver"
        ]
        author = "Patrick Pacher"
        license = "MIT"
        description = "Add support to use DNSCrypt servers"
        artifact_template = "{{source}}/releases/download/{{version}}/{{archive_file}}_{{stripped_version}}_{{os}}_{{arch}}.tar.gz"
        archive_file = "portmaster-plugin-dnscrypt"
}