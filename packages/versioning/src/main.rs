use anyhow::Error;
use colored::Colorize;
use inquire::{required, Text};
use macros_rs::{crashln, str};
use regex::Regex;
use serde::Deserialize;
use serde_json::{from_str, to_string_pretty, Value};
use std::io::prelude::*;
use std::process::Command;
use std::str::from_utf8;
use std::{fs, fs::File};
use toml_edit::{value, Document};

const MAIDFILE: &str = "Maidfile.toml";
const VERSION_PKG: &str = "packages/versioning/Cargo.toml";
const DATABASE_PKG: &str = "packages/database/Cargo.toml";
const AUTHENTICATION_PKG: &str = "packages/authentication/config/version.go";
const SERVER_PKG: &str = "packages/server/package.json";

struct Files {
    maidfile: String,
    versioning: String,
    database: String,
    authentication: String,
    server: String,
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
        versioning: match fs::read_to_string(VERSION_PKG) {
            Ok(file) => file,
            Err(err) => crashln!("{err}"),
        },
        database: match fs::read_to_string(DATABASE_PKG) {
            Ok(file) => file,
            Err(err) => crashln!("{err}"),
        },
        authentication: match fs::read_to_string(AUTHENTICATION_PKG) {
            Ok(file) => file,
            Err(err) => crashln!("{err}"),
        },
        server: match fs::read_to_string(SERVER_PKG) {
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

    let regex = Regex::new(r#"const version = "([0-9]{1,4}(\.[0-9a-z]{1,6}){1,5})""#).unwrap();
    let authentication = regex.replace_all(str!(read().authentication.clone()), format!("const version = \"{}\"", version));

    let mut maidfile: Document = match read().maidfile.parse() {
        Ok(doc) => doc,
        Err(err) => crashln!("{err}"),
    };

    let mut versioning: Document = match read().versioning.parse() {
        Ok(doc) => doc,
        Err(err) => crashln!("{err}"),
    };

    let mut database: Document = match read().database.parse() {
        Ok(doc) => doc,
        Err(err) => crashln!("{err}"),
    };

    let mut server: Value = match from_str(&read().server) {
        Ok(doc) => doc,
        Err(err) => crashln!("{err}"),
    };

    maidfile["project"]["version"] = value(version.clone());
    versioning["package"]["version"] = value(version.clone());
    database["package"]["version"] = value(version.clone());
    server["version"] = Value::String(version.clone());

    write(MAIDFILE, maidfile.to_string());
    write(VERSION_PKG, versioning.to_string());
    write(DATABASE_PKG, database.to_string());
    write(AUTHENTICATION_PKG, authentication.to_string());
    write(SERVER_PKG, to_string_pretty(&server).unwrap());
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
