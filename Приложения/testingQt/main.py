import sys
from PyQt6.QtWidgets import (QApplication, QMainWindow, QTableView, 
                             QVBoxLayout, QWidget, QLabel, QMessageBox)
from PyQt6.QtSql import QSqlDatabase, QSqlTableModel
from PyQt6.QtCore import Qt

class MultiDbApp(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("PyQt6: Связь PostgreSQL + MySQL")
        self.resize(600, 500)

        # 1. Инициализация подключений
        if not self.init_connections():
            sys.exit(1)

        # 2. Создание интерфейса
        layout = QVBoxLayout()

        # --- Верхняя часть (PostgreSQL) ---
        layout.addWidget(QLabel("Группы (PostgreSQL):"))
        self.master_view = QTableView()
        self.master_model = QSqlTableModel(db=QSqlDatabase.database("pg_conn"))
        self.master_model.setTable("groups")
        self.master_model.select()
        self.master_view.setModel(self.master_model)
        layout.addWidget(self.master_view)

        # --- Нижняя часть (MySQL) ---
        layout.addWidget(QLabel("Студенты (MySQL):"))
        self.detail_view = QTableView()
        self.detail_model = QSqlTableModel(db=QSqlDatabase.database("mysql_conn"))
        self.detail_model.setTable("students")
        self.detail_model.select()
        self.detail_view.setModel(self.detail_model)
        layout.addWidget(self.detail_view)

        # 3. Связка (Логика Master-Detail)
        self.master_view.clicked.connect(self.on_group_selected)

        # Рендеринг
        container = QWidget()
        container.setLayout(layout)
        self.setCentralWidget(container)

    def init_connections(self):
        # Подключаем PG
        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName("127.0.0.1")
        db_pg.setDatabaseName("master_detail")
        db_pg.setUserName("postgres")
        db_pg.setPassword("1234")

        if not db_pg.open():
            QMessageBox.critical(None, "Ошибка PG", db_pg.lastError().text())
            return False

        # Подключаем MySQL
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn")
        db_ms.setHostName("127.0.0.1")
        db_ms.setDatabaseName("master_detail")
        db_ms.setUserName("root")
        db_ms.setPassword("1234")

        if not db_ms.open():
            QMessageBox.critical(None, "Ошибка MySQL", db_ms.lastError().text())
            return False

        return True

    def on_group_selected(self, index):
        """ Срабатывает при клике на строку в верхней таблице """
        # Получаем ID выбранной группы (колонка 0)
        group_id = self.master_model.data(self.master_model.index(index.row(), 0))
        
        # Фильтруем нижнюю таблицу
        self.detail_model.setFilter(f"group_id = {group_id}")
        self.detail_model.select()

if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = MultiDbApp()
    window.show()
    sys.exit(app.exec())
