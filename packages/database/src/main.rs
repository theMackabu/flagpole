#[macro_use]
extern crate actix_web;
extern crate clap;

use actix_web::http::StatusCode;
use actix_web::{middleware, web, App, HttpResponse, HttpServer};
use serde::{Deserialize, Serialize};
use serde_json::json;
use sled::Db;
use std::env;
use std::sync::{Arc, Mutex};

use tracing_actix_web::TracingLogger;
use tracing_bunyan_formatter::{BunyanFormattingLayer, JsonStorageLayer};
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::{EnvFilter, Registry};

#[derive(Serialize, Deserialize)]
struct DbConfig {
    name: String,
    path: String,
}

struct ServerState {
    name: String,
    db: Db,
}

fn err_not_found() -> HttpResponse {
    HttpResponse::build(StatusCode::NOT_FOUND).content_type("application/json").body(
        json!({
          "error": {
             "code" : -404,
              "message": "not found"}})
        .to_string(),
    )
}

fn err_500() -> HttpResponse {
    HttpResponse::build(StatusCode::INTERNAL_SERVER_ERROR).content_type("application/json").body(
        json!({
          "error": {
             "code" : -500,
              "message": "internal server error"}})
        .to_string(),
    )
}

fn ok_binary(val: Vec<u8>) -> HttpResponse {
    HttpResponse::Ok().content_type("application/octet-stream").body(val)
}

fn ok_json(jval: serde_json::Value) -> HttpResponse {
    HttpResponse::Ok().content_type("application/json").body(jval.to_string())
}

#[get("/")]
async fn req_index(m_state: web::Data<Arc<Mutex<ServerState>>>) -> HttpResponse {
    let state = m_state.lock().unwrap();

    ok_json(json!({
        "name": "database",
        "version": env!("CARGO_PKG_VERSION"),
        "databases": [
            { "name": state.name }
        ]
    }))
}

async fn req_delete(m_state: web::Data<Arc<Mutex<ServerState>>>, path: web::Path<(String, String)>) -> HttpResponse {
    let state = m_state.lock().unwrap();

    if state.name != path.0 {
        return err_not_found();
    }

    match state.db.remove(path.1.clone()) {
        Ok(optval) => match optval {
            Some(_val) => ok_json(json!({"result": true})),
            None => err_not_found(),
        },
        Err(_e) => err_500(),
    }
}

async fn req_get(m_state: web::Data<Arc<Mutex<ServerState>>>, path: web::Path<(String, String)>) -> HttpResponse {
    let state = m_state.lock().unwrap();

    if state.name != path.0 {
        return err_not_found();
    }

    match state.db.get(path.1.clone()) {
        Ok(optval) => match optval {
            Some(val) => ok_binary(val.to_vec()),
            None => err_not_found(),
        },
        Err(_e) => err_500(),
    }
}

async fn req_put(m_state: web::Data<Arc<Mutex<ServerState>>>, (path, body): (web::Path<(String, String)>, web::Bytes)) -> HttpResponse {
    let state = m_state.lock().unwrap();

    if state.name != path.0 {
        return err_not_found();
    }

    match state.db.insert(path.1.as_str(), body.to_vec()) {
        Ok(_optval) => ok_json(json!({"result": true})),
        Err(_e) => err_500(),
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();

    let env_filter = EnvFilter::try_from_default_env().unwrap_or(EnvFilter::new("info"));
    let formatting_layer = BunyanFormattingLayer::new("database".to_string(), std::io::stdout);
    let subscriber = Registry::default().with(env_filter).with(JsonStorageLayer).with(formatting_layer);
    tracing::subscriber::set_global_default(subscriber).expect("Failed to install `tracing` subscriber.");

    let matches = clap::Command::new("database")
        .arg(clap::Arg::new("address").long("address").help("Server address").required(true))
        .arg(clap::Arg::new("port").long("port").help("Server port").required(true))
        .arg(clap::Arg::new("database").long("database").help("Database Info").required(true))
        .get_matches();

    let address = matches.get_one::<String>("address").unwrap();
    let port = matches.get_one::<String>("port").unwrap();
    let database = matches.get_one::<String>("database").unwrap();

    let pair = format!("{}:{}", address, port);
    let server_hdr = format!("database/{}", env!("CARGO_PKG_VERSION"));

    let db_config = sled::Config::default().path(database).use_compression(false);
    let db = db_config.open().unwrap();

    let srv_state = Arc::new(Mutex::new(ServerState {
        name: database.clone(),
        db: db.clone(),
    }));

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(Arc::clone(&srv_state)))
            .wrap(middleware::DefaultHeaders::new().add(("Server", server_hdr.to_string())))
            .wrap(TracingLogger::default())
            .service(req_index)
            .service(
                web::resource("/api/{db}/{key}")
                    .route(web::get().to(req_get))
                    .route(web::put().to(req_put))
                    .route(web::delete().to(req_delete)),
            )
    })
    .bind(pair.to_string())?
    .run()
    .await
}
