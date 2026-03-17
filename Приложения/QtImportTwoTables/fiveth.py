import sys
from PyQt6.QtWidgets import (QApplication, QMainWindow, QVBoxLayout, QHBoxLayout, 
                             QWidget, QLabel, QMessageBox, QSplitter, 
                             QComboBox, QScrollArea, QFrame)
from PyQt6.QtSql import QSqlDatabase, QSqlQueryModel, QSqlQuery
from PyQt6.QtCore import Qt
from PyQt6.QtGui import QFont

# Настройки подключения
MYSQL_HOST, MYSQL_DB, MYSQL_USER, MYSQL_PASS = "127.0.0.1", "insurance", "root", "1234"
PG_HOST, PG_DB, PG_USER, PG_PASS = "127.0.0.1", "insurance", "postgres", "1234"

# SQL Запросы
SQL_GET_EMPLOYEES = "SELECT employee_id, full_name FROM employees"
SQL_GET_POLICYHOLDERS = "SELECT policy_number, full_name FROM policyholders WHERE employee_id = :emp_id"
SQL_GET_CLAIMS = "SELECT claim_id, description, event_date, payout FROM claims WHERE policy_number = :policy_num"

# --- АНАЛОГ USERCONTROL (Карточка с отображением ключей и данных) ---
class ClaimCard(QFrame):
    def __init__(self, emp_id, policy_num, claim_id, description, date, payout):
        super().__init__()
        self.setFrameShape(QFrame.Shape.StyledPanel)
        self.setStyleSheet("""
            ClaimCard {
                background-color: #fcfcfc;
                border: 1px solid #d1d1d1;
                border-radius: 8px;
                margin-bottom: 10px;
            }
            QLabel#KeyLabel {
                color: #7f8c8d;
                font-size: 10px;
                text-transform: uppercase;
            }
            QLabel#KeyValue {
                color: #e67e22;
                font-weight: bold;
                font-family: 'Consolas';
            }
            QLabel#DataText {
                font-size: 14px;
                color: #2c3e50;
            }
            QLabel#Payout {
                font-size: 16px;
                font-weight: bold;
                color: #27ae60;
            }
        """)

        main_layout = QVBoxLayout()

        # 1. Секция ключей (Путь, который привел к записи)
        keys_layout = QHBoxLayout()
        
        def add_key(label, value):
            v_box = QVBoxLayout()
            lbl = QLabel(label); lbl.setObjectName("KeyLabel")
            val = QLabel(str(value)); val.setObjectName("KeyValue")
            v_box.addWidget(lbl)
            v_box.addWidget(val)
            keys_layout.addLayout(v_box)

        add_key("Emp ID", emp_id)
        keys_layout.addSpacing(20)
        add_key("Policy №", policy_num)
        keys_layout.addSpacing(20)
        add_key("Claim ID", claim_id)
        keys_layout.addStretch()
        
        main_layout.addLayout(keys_layout)
        
        # Разделительная линия
        line = QFrame(); line.setFrameShape(QFrame.Shape.HLine); line.setFrameShadow(QFrame.Shadow.Sunken)
        main_layout.addWidget(line)

        # 2. Секция данных
        lbl_desc = QLabel(f"Описание: {description}")
        lbl_desc.setObjectName("DataText")
        lbl_desc.setWordWrap(True)
        
        lbl_date = QLabel(f"Дата происшествия: {date}")
        lbl_date.setObjectName("DataText")

        lbl_payout = QLabel(f"Сумма выплаты: {payout} ₽")
        lbl_payout.setObjectName("Payout")
        lbl_payout.setAlignment(Qt.AlignmentFlag.AlignRight)

        main_layout.addWidget(lbl_desc)
        main_layout.addWidget(lbl_date)
        main_layout.addWidget(lbl_payout)

        self.setLayout(main_layout)

# --- ОСНОВНАЯ СЕКЦИЯ БД ---
class DatabaseSection(QWidget):
    def __init__(self, title, db_connection_name):
        super().__init__()
        self.db_name = db_connection_name
        self.current_emp_id = None
        self.init_ui(title)
        self.load_employees()

    def init_ui(self, title):
        layout = QVBoxLayout()
        
        # Заголовок движка
        lbl_title = QLabel(title)
        lbl_title.setFont(QFont("Bahnschrift", 22))
        lbl_title.setStyleSheet("color: #2980b9;")
        layout.addWidget(lbl_title)

        # Комбобоксы
        grid_layout = QHBoxLayout()
        
        v_emp = QVBoxLayout()
        v_emp.addWidget(QLabel("1. Выберите сотрудника:"))
        self.combo_emp = QComboBox()
        self.combo_emp.currentIndexChanged.connect(self.on_employee_changed)
        v_emp.addWidget(self.combo_emp)
        
        v_hold = QVBoxLayout()
        v_hold.addWidget(QLabel("2. Выберите клиента:"))
        self.combo_holder = QComboBox()
        self.combo_holder.currentIndexChanged.connect(self.on_policyholder_changed)
        v_hold.addWidget(self.combo_holder)
        
        grid_layout.addLayout(v_emp)
        grid_layout.addLayout(v_hold)
        layout.addLayout(grid_layout)

        # Область для карточек (UserControls)
        layout.addWidget(QLabel("Результаты (Claims):"))
        self.scroll = QScrollArea()
        self.scroll.setWidgetResizable(True)
        self.scroll_content = QWidget()
        self.cards_layout = QVBoxLayout(self.scroll_content)
        self.cards_layout.setAlignment(Qt.AlignmentFlag.AlignTop)
        self.scroll.setWidget(self.scroll_content)
        layout.addWidget(self.scroll)

        self.setLayout(layout)

        # Модели
        self.model_emp = QSqlQueryModel()
        self.model_holder = QSqlQueryModel()
        self.combo_emp.setModel(self.model_emp)
        self.combo_holder.setModel(self.model_holder)

    def load_employees(self):
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.exec(SQL_GET_EMPLOYEES)
        self.model_emp.setQuery(query)
        self.combo_emp.setModelColumn(1)

    def on_employee_changed(self, index):
        if index < 0: return
        self.current_emp_id = self.model_emp.record(index).value("employee_id")
        
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.prepare(SQL_GET_POLICYHOLDERS)
        query.bindValue(":emp_id", self.current_emp_id)
        query.exec()
        
        self.model_holder.setQuery(query)
        self.combo_holder.setModelColumn(1)

    def on_policyholder_changed(self, index):
        # Очистка старых карточек
        while self.cards_layout.count():
            item = self.cards_layout.takeAt(0)
            if item.widget(): item.widget().deleteLater()

        if index < 0: return
        
        policy_num = self.model_holder.record(index).value("policy_number")
        
        db = QSqlDatabase.database(self.db_name)
        query = QSqlQuery(db)
        query.prepare(SQL_GET_CLAIMS)
        query.bindValue(":policy_num", policy_num)
        query.exec()

        # Создаем карточки (UserControls)
        while query.next():
            card = ClaimCard(
                emp_id=self.current_emp_id,
                policy_num=policy_num,
                claim_id=query.value("claim_id"),
                description=query.value("description"),
                date=query.value("event_date"),
                payout=query.value("payout")
            )
            self.cards_layout.addWidget(card)

# --- ГЛАВНОЕ ОКНО ---
class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Insurance Key-Value Cards")
        self.resize(1200, 800)

        if not self.init_databases():
            sys.exit(1)

        splitter = QSplitter(Qt.Orientation.Horizontal)
        splitter.addWidget(DatabaseSection("MySQL", "mysql_conn"))
        splitter.addWidget(DatabaseSection("PostgreSQL", "pg_conn"))
        
        self.setCentralWidget(splitter)

    def init_databases(self):
        db_ms = QSqlDatabase.addDatabase("QMARIADB", "mysql_conn")
        db_ms.setHostName(MYSQL_HOST); db_ms.setDatabaseName(MYSQL_DB)
        db_ms.setUserName(MYSQL_USER); db_ms.setPassword(MYSQL_PASS)

        db_pg = QSqlDatabase.addDatabase("QPSQL", "pg_conn")
        db_pg.setHostName(PG_HOST); db_pg.setDatabaseName(PG_DB)
        db_pg.setUserName(PG_USER); db_pg.setPassword(PG_PASS)

        if not db_ms.open():
            QMessageBox.critical(self, "Error", f"MySQL Error: {db_ms.lastError().text()}")
            return False
        if not db_pg.open():
            QMessageBox.critical(self, "Error", f"PostgreSQL Error: {db_pg.lastError().text()}")
            return False
        return True

if __name__ == "__main__":
    app = QApplication(sys.argv)
    app.setFont(QFont("Segoe UI", 10))
    window = MainWindow()
    window.show()
    sys.exit(app.exec())