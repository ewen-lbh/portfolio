#!/usr/bin/env ruby
require "uri"
require "net/http"
require "json"
require "set"
require "yaml"

require "bundler/inline"
gemfile true do
  source 'https://rubygems.org' 
  gem 'parallel'
end

# dumb websites
# maps hosts to allowed >=400 http status codes
exceptions = {
  "twitter.com": 400,
  "www.fiverr.com": 403,
  "distrokid.com": 403,
  "www.amazon.com": 500,
  "evann.tech": 0,
}


check_every = 3600 / 8 # seconds


def log_error msg
  File.write "deadlinks.log", "[#{Time.now}] #{msg}\n", mode: 'a+'
end

def parse_uri uri
  if uri !~ /https?:\/\// then uri = "https://ewen.works#{uri}" end
  return URI.parse URI::Parser.new.escape uri
end

def check_broken_links database, exceptions 
  links = Set.new

  database.each do |id, work|
    work['content'].each do |language, localized|
      localized['blocks'].each do |block|
        if block['type'] == 'link' 
          links.add? parse_uri block['url']
        end
      end
    end
  end

  YAML.load_file("./tags.yaml").each do |tag|
    if tag.has_key? "learn more at"
      links.add? parse_uri tag["learn more at"]
    end
  end

  YAML.load_file("./technologies.yaml").each do |tag|
    if tag.has_key? "learn more at"
      links.add? parse_uri tag["learn more at"]
    end
  end

  return Parallel.map(links) do |url|
    begin
      puts "Checking #{url}"
      response = (Net::HTTP.get_response url)
      if response.code.to_i >= 400
        unless exceptions.has_key? url.host.to_sym and response.code.to_i == exceptions[url.host.to_sym]
          log_error "#{url} broken: got status #{response.code}"
          url
        end
      end
    rescue StandardError => ex
      unless exceptions.has_key? url.host.to_sym and exceptions[url.host.to_sym] == 0
        log_error "#{url}: utterly broken: #{ex}"
        url
      end
    end
  end.compact
end

def push_kuma_status error_messages
  uptime_kuma_push_url = URI.parse ARGV[1]
  push_params = {
    "status": if error_messages.empty? then "up" else "down" end,
      "msg": if error_messages.empty? then "OK" else error_messages.join " \n" end,
      "ping": ""
  }

  uptime_kuma_push_url.query = URI.encode_www_form push_params.to_a

  Net::HTTP.get_response uptime_kuma_push_url
end

while true
  database = JSON.load (File.open ARGV[0])
  broken_links = check_broken_links database, exceptions
  push_kuma_status broken_links
  puts "Will check again in #{check_every} seconds…"
  sleep check_every
end
