#!/usr/bin/env ruby

if ["-h", "--help"].include? ARGV[0] 
  puts "Usage: #{__FILE__}"
  puts "              Include all recipes"
  puts "       #{__FILE__} including recipe1 recipe2 ..."
  puts "              Only include the specified recipes"
  puts "       #{__FILE__} excluding recipe1 recipe2 ..."
  puts "              Exclude the specified recipes"
  exit 0
end

def titlecase(s)
  s = s.dup
  s.sub! /[-_]/, " "
  s[0].upcase + s[1..-1]
end

# get Justfile contents

justfile = File.read("Justfile")

# parse rudimentarily into a hash of recipe => array of commands

recipes = {}
current_recipe = nil
indentation = "    "


def remove_args(recipe_name_line)
  recipe_name_line.split(" ").first
end

def render_args_substitutions(commandline)
  commandline.sub(/\{\{\s*([^\s}]+)\s*\}\}/) { |match| "YOUR_#{$1.upcase}" }
end

def render_just_calls commandline
  unless commandline.start_with? "just"
    return commandline
  end
  args = commandline.split(" ")[1..-1]
  "# Run #{titlecase args.first}'s commands"
end

def render_commandline commandline
  render_just_calls render_args_substitutions (commandline.strip)
end

justfile.each_line do |line|
  recipe_name = remove_args line.strip.sub /:\s*$/, ""
  if current_recipe == nil
    current_recipe = recipe_name
  elsif line.start_with? indentation
    if recipes.has_key? current_recipe
      recipes[current_recipe].push render_commandline line
    else
      $stderr.puts "Creating new recipe #{current_recipe}"
      recipes[current_recipe] = [render_commandline(line)]
    end
    $stderr.puts "Added #{line.strip} to #{current_recipe}"
  else
    current_recipe = recipe_name
  end
end

$stderr.puts "Parsed Justfile: #{recipes}"



def recipe_markdown(recipe, commands)
  [
    "## #{titlecase recipe}",
    "",
    "```sh",
    *commands,
    "```",
  ].join "\n"
end

excluded_recipes = []
if ARGV[0] == "excluding"
  excluded_recipes = ARGV[1..-1]
elsif ARGV[0] == "including"
  excluded_recipes = recipes.keys - ARGV[1..-1]
end

rendered = recipes.reject { |name| excluded_recipes.include? name }.map do |name, commands|
  unless commands.empty?
    recipe_markdown name, commands
  end
end.compact.join "\n\n"


puts rendered
