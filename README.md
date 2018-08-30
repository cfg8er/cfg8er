# Cfg8er

## Goals

- Service where the input configuration data is the source of truth. The input is the Git repo. No intermediate translation or storage layer.
- Does not present the traditional filesystem checkout of the repo. Think turning a Git repo directly into an API.
- Looks up Git objects in a repo based on path and version and serves it as configuration data.
- Common semantic versioning based interface for configuration data, eg. /r/example_repo/v1/path/config.yml will present path/config.yml from the latest v1.x.y tag.
- Interface will be HTTP or filesystem to start with. 
- Event sourced data store since it's using Git.
- Allows for common Git and pull request workflow for admins, operators, developers.
- Promote concept single-source of truth.
- Provide configurable merging behavior across directory structure.
- Provide limited useful macros that directly contribute to single source of truth but avoids magic.

### Ideas for Future

- Compatible with Ansible host_var + group_var, Puppet hiera, iPXE scripts
- Compatibility interface for Warewulf's http interface, which converts YAML (or other formats) into a flattened shell variable format that Warewulf's initrd is expecting.
