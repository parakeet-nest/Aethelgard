use extism_pdk::*;

#[derive(serde::Deserialize)]
struct Arguments {
    from: u32,
    to: u32,
}

fn can_move(from: u32, to: u32) -> bool {
    // Define the connections for each room as a list of tuples (room number, list of connected rooms)
    let connections = vec![
        (1, vec![2, 6]),
        (2, vec![1, 3, 7]),
        (3, vec![2, 4, 8]),
        (4, vec![3, 5, 9]),
        (5, vec![4, 10]),
        (6, vec![1, 7, 11]),
        (7, vec![2, 6, 8, 12]),
        (8, vec![3, 7, 9, 13]),
        (9, vec![4, 8, 10, 14]),
        (10, vec![5, 9, 15]),
        (11, vec![6, 12, 16]),
        (12, vec![7, 11, 13, 17]),
        (13, vec![8, 12, 14, 18]),
        (14, vec![9, 13, 15, 19]),
        (15, vec![10, 14, 20]),
        (16, vec![11, 17, 21]),
        (17, vec![12, 16, 18, 22]),
        (18, vec![13, 17, 19, 23]),
        (19, vec![14, 18, 20, 24]),
        (20, vec![15, 19, 25]),
        (21, vec![16, 22]),
        (22, vec![17, 21, 23]),
        (23, vec![18, 22, 24]),
        (24, vec![19, 23, 25]),
        (25, vec![20, 24]),
    ];

    // Iterate through the connections and check if the "from" room has a connection to the "to" room
    for (room, connected_rooms) in connections {
        if room == from {
            return connected_rooms.contains(&to);
        }
    }

    // If no connection is found, return false
    false
}


#[plugin_fn]
pub fn move_person(Json(args): Json<Arguments>) -> FnResult<String> {

    let can_i_move: bool = can_move(args.from, args.to);

    if can_i_move {
        Ok(format!("🙂 🦶 from {} to {}", args.from, args.to))
    } else {
        Ok(format!("😡 ✋ from {} to {}", args.from, args.to))
    }

    
}
// rustup target add wasm32-unknown-unknown
// cargo build --target wasm32-unknown-unknown
