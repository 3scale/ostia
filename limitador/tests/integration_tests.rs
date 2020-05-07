extern crate limitador;

use limitador::limit::Limit;
use limitador::RateLimiter;
use std::collections::{HashMap, HashSet};

#[test]
fn add_a_limit() {
    let limit = Limit::new(
        "test_namespace",
        10,
        60,
        vec!["req.method == GET"],
        vec!["req.method", "app_id"],
    );

    let mut rate_limiter = RateLimiter::new();
    rate_limiter.add_limit(limit.clone());

    let mut expected_result = HashSet::new();
    expected_result.insert(limit);

    assert_eq!(rate_limiter.get_limits("test_namespace"), expected_result)
}

// TODO: test add multiple limits same namespace.

#[test]
fn rate_limited() {
    let max_hits = 3;

    let limit = Limit::new(
        "test_namespace",
        max_hits,
        60,
        vec!["req.method == GET"],
        vec!["req.method", "app_id"],
    );

    let mut rate_limiter = RateLimiter::new();
    rate_limiter.add_limit(limit.clone());

    let mut values: HashMap<String, String> = HashMap::new();
    values.insert("namespace".to_string(), "test_namespace".to_string());
    values.insert("req.method".to_string(), "GET".to_string());
    values.insert("app_id".to_string(), "test_app_id".to_string());

    for _ in 0..max_hits {
        assert_eq!(false, rate_limiter.is_rate_limited(&values).unwrap());
        rate_limiter.update_counters(&values).unwrap();
    }
    assert_eq!(true, rate_limiter.is_rate_limited(&values).unwrap());
}