use anyhow::Error;
use colored::Colorize;
use inquire::{required, Text};
use macros_rs::{crashln, str};
use serde::Deserialize;
use std::io::prelude::*;
use std::process::Command;
use std::str::from_utf8;
use std::{fs, fs::File};
use toml_edit::{value, Document};

const MAIDFILE: &str = "Maidfile.toml";
const VERSION_TOML: &str = "packages/versioning/Cargo.toml";
const DATABASE_TOML: &str = "packages/database/Cargo.toml";

struct Files {
    maidfile: String,
    version_system: String,
    database: String,
}

#[derive(Deserialize)]
struct Base {
    project: Project,
}

#[derive(Deserialize)]
struct Project {
    version: String,
}

fn version() -> Result<String, Error> {
    let data = Command::new("maid").args(["butler", "json"]).output()?;
    let output = from_utf8(&data.stdout)?;
    let json: Base = serde_json::from_str(output)?;

    Ok(json.project.version)
}

fn read() -> Files {
    return Files {
        maidfile: match fs::read_to_string(MAIDFILE) {
            Ok(file) => file,
            Err(err) => crashln!("{err}"),
        },
        version_system: match fs::read_to_string(VERSION_TOML) {
            Ok(file) => file,
            Err(err) => crashln!("{err}"),
        },
        database: match fs::read_to_string(DATABASE_TOML) {
            Ok(file) => file,
            Err(err) => crashln!("{err}"),
        },
    };
}

fn write(path: &str, contents: String) {
    let mut file = match File::create(path) {
        Ok(file) => file,
        Err(err) => crashln!("{err}"),
    };

    match file.write_all(contents.as_bytes()) {
        Ok(file) => file,
        Err(err) => crashln!("{err}"),
    };

    println!(" {} {}", "-".white(), format!("{path}").yellow());
}

fn set(version: String) {
    println!("\n{} {} {}", "Updated versions to".white(), format!("v{version}").bright_green(), "in:".white());

    let mut maidfile = match read().maidfile.parse::<Document>() {
        Ok(doc) => doc,
        Err(err) => crashln!("{err}"),
    };

    let mut version_system = match read().version_system.parse::<Document>() {
        Ok(doc) => doc,
        Err(err) => crashln!("{err}"),
    };

    let mut database = match read().database.parse::<Document>() {
        Ok(doc) => doc,
        Err(err) => crashln!("{err}"),
    };

    maidfile["project"]["version"] = value(version.clone());
    version_system["package"]["version"] = value(version.clone());
    database["package"]["version"] = value(version.clone());

    write(MAIDFILE, maidfile.to_string());
    write(VERSION_TOML, version_system.to_string());
    write(DATABASE_TOML, database.to_string());
}

fn main() {
    let version = match version() {
        Ok(version) => version,
        Err(err) => crashln!("{err}"),
    };

    match Text::new("Package version?").with_initial_value(str!(version)).with_validator(required!()).prompt() {
        Ok(data) => set(data),
        Err(_) => crashln!("Aborting..."),
    }
}
