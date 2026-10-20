from __future__ import annotations

from db_manager import DatabaseError, DatabaseManager
from PyQt6.QtWidgets import (
    QDialog,
    QDialogButtonBox,
    QFormLayout,
    QLineEdit,
    QMessageBox,
    QSpinBox,
)


class ConnectionDialog(QDialog):
    """Окно ввода параметров подключения к MySQL при старте приложения."""

    def __init__(self, parent=None):
        super().__init__(parent)
        self.setWindowTitle("Подключение к базе данных")
        self.db_manager = DatabaseManager()

        self.host_edit = QLineEdit("127.0.0.1")
        self.port_spin = QSpinBox()
        self.port_spin.setRange(1, 65535)
        self.port_spin.setValue(3306)
        self.user_edit = QLineEdit("root")
        self.password_edit = QLineEdit()
        self.password_edit.setEchoMode(QLineEdit.EchoMode.Password)
        self.database_edit = QLineEdit("insurance")

        layout = QFormLayout(self)
        layout.addRow("Хост:", self.host_edit)
        layout.addRow("Порт:", self.port_spin)
        layout.addRow("Пользователь:", self.user_edit)
        layout.addRow("Пароль:", self.password_edit)
        layout.addRow("База данных:", self.database_edit)

        buttons = QDialogButtonBox(
            QDialogButtonBox.StandardButton.Ok | QDialogButtonBox.StandardButton.Cancel
        )
        buttons.accepted.connect(self._on_ok)
        buttons.rejected.connect(self.reject)
        layout.addRow(buttons)

    def _on_ok(self) -> None:
        try:
            self.db_manager.connect(
                host=self.host_edit.text().strip(),
                port=self.port_spin.value(),
                user=self.user_edit.text().strip(),
                password=self.password_edit.text(),
                database=self.database_edit.text().strip(),
            )
        except DatabaseError as exc:
            QMessageBox.critical(self, "Ошибка подключения", str(exc))
            return
        self.accept()
