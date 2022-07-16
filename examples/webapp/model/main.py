import logging


# from model.initial import Initializer

def start_predict(data: list) -> int:
    """
    Пример cpu-bounded функции
    :return:
    """
    # print(Initializer.DF_SHARED_VALUE.df)
    try:
        sum_predict = 0
        for num in data:
            sum_predict += sum(list(fibb_gen(num)))

        res = len(str(2 ** sum_predict))

        logging.info('DONE')
        logging.info(res)
        return res
    except Exception as e:
        logging.info(e)


def fibb_gen(n):
    a, b = 0, 1
    for _ in range(n):
        yield a
        a, b = b, a + b
