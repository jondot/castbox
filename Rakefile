
include FileUtils

# 
# fpm?
# fpm -s tar -t osxpkg --osxpkg-ownership preserve --prefix /Applications/Castbox.app castbox-darwin-amd64.tar.gz
#

ARCHS = [
  ['windows-386', ''],
  ['linux-arm', ''],
  ['linux-amd64',''],
  ['linux-386',''],
  ['darwin-amd64','/Applications/Castbox.app']
]

task :package do
  ARCHS.each do |arch, prefix|
    `mkdir build`
    `mkdir -p pkg/#{arch}`
    bin = go_build("build/castbox", arch)

    cd("pkg/#{arch}") do
      cp "../../Castfile", "Castfile"
      cp "../../#{bin}", "castbox"

      tarball = "castbox-#{arch}.tar.gz"
      `tar -czf #{tarball} *`
      cp tarball, "../"
    end

    `rm -rf pkg/#{arch}`
    `rm -rf build`
  end
end


def go_build(label, arch)
    system({"GOOS" => arch.split('-')[0], 
          "GOARCH" => arch.split('-')[1]}, 
          "go build -o #{label}-#{arch}")
    "#{label}-#{arch}"
end

