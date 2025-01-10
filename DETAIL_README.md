## Entities

- Inventory
- UI :: WaveController, Inventory
    - Hat :: Inventory(To add stuff), WaveController(to modify next wave risk)
- Grid (manages tower entity list)
    - Tower :: WaveController(To check for enemy positions for shooting)
- WaveController (manages entities) :: Grid(to check collision/path)
    - Enemy
