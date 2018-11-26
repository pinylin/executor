// app.rs
use actix::prelude::*;
use actix_web::{http::Method, middleware, App};
use models::DbExecutor;

pub struct AppState {
    pub db: Addr<DbExecutor>,
}

// helper function to create and returns the app after mounting all routes/resources
pub fn create_app(db: Addr<DbExecutor>) -> App<AppState> {
    App::with_state(AppState { db })
        // setup builtin logger to get nice logging for each request
        .middleware(middleware::Logger::new("\"%r\" %s %b %Dms"))

        // routes for authentication
        .resource("/auth", |r| {
        })
        // routes to invitation
        .resource("/invitation/", |r| {
        })
        // routes to register as a user after the
        .resource("/register/", |r| {
        })
}