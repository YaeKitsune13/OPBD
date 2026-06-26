import sys

from PyQt6.QtWidgets import QApplication, QDialog

from connection_dialog import ConnectionDialog
from main_window import MainWindow


def main() -> None:
    app = QApplication(sys.argv)

    conn_dialog = ConnectionDialog()
    if conn_dialog.exec() != QDialog.DialogCode.Accepted:
        sys.exit(0)

    window = MainWindow(conn_dialog.db_manager)
    window.show()

    exit_code = app.exec()
    conn_dialog.db_manager.close()
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
