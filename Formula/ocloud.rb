# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Ocloud < Formula
  desc "Tool for finding and connecting to OCI instances"
  homepage "https://github.com/rozdolsky33/ocloud"
  version "0.0.10"
  license "MIT"

  on_macos do
    url "https://github.com/rozdolsky33/ocloud/releases/download/v0.0.10/ocloud_0.0.10_darwin_all.tar.gz"
    sha256 "f0c6eefce56b3aa54ff2e3b4cdfef46fd3b8cfe90bf7097322c0a4c3d486bb36"

    def install
      bin.install "ocloud"
    end
  end

  on_linux do
    if Hardware::CPU.intel? and Hardware::CPU.is_64_bit?
      url "https://github.com/rozdolsky33/ocloud/releases/download/v0.0.10/ocloud_0.0.10_linux_amd64.tar.gz"
      sha256 "051b9f410dcab7f5f6318c0bdd870d9b2816d2b0ad49a56178111cdb4dd4f5de"
      def install
        bin.install "ocloud"
      end
    end
    if Hardware::CPU.arm? and Hardware::CPU.is_64_bit?
      url "https://github.com/rozdolsky33/ocloud/releases/download/v0.0.10/ocloud_0.0.10_linux_arm64.tar.gz"
      sha256 "b108485adbb4acea60aa3e1a165821dd483de990e2779d3742aa81c4fc8b29c9"
      def install
        bin.install "ocloud"
      end
    end
  end
end
